package handler

import (
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"ecbms/internal/model"

	"github.com/gin-gonic/gin"
)

func GetNetworkInterfaces(c *gin.Context) {
	ifaces, err := net.Interfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network interfaces"})
		return
	}

	var interfaces []model.NetworkInterface
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			ni := model.NetworkInterface{
				Name:   iface.Name,
				HWAddr: iface.HardwareAddr.String(),
				Status: "up",
			}

			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.To4() != nil {
					}
				}
			}

			interfaces = append(interfaces, ni)
		}
	}

	c.JSON(http.StatusOK, interfaces)
}

func GetNetworkConfigs(c *gin.Context) {
	configs, err := model.GetAllNetworkConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network configs"})
		return
	}
	c.JSON(http.StatusOK, configs)
}

type NetworkConfigRequest struct {
	InterfaceName string `json:"interface_name" binding:"required"`
	Mode          string `json:"mode" binding:"required,oneof=dhcp static"`
	IPAddress     string `json:"ip_address"`
	Netmask       string `json:"netmask"`
	Gateway       string `json:"gateway"`
	DNSServers    string `json:"dns_servers"`
}

func ConfigureNetwork(c *gin.Context) {
	var req NetworkConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Mode == "static" {
		if req.IPAddress == "" || req.Netmask == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IP address and netmask are required for static mode"})
			return
		}
	}

	if err := model.SaveNetworkConfig(
		req.InterfaceName, req.Mode, req.IPAddress, req.Netmask, req.Gateway, req.DNSServers,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save network config"})
		return
	}

	username, _ := c.Get("username")
	model.AddSystemLog("network", "Network config updated for "+req.InterfaceName+" by "+username.(string), "info")

	c.JSON(http.StatusOK, gin.H{"message": "Network configuration saved successfully"})
}

func ApplyNetworkConfig(c *gin.Context) {
	if runtime.GOOS != "linux" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Network configuration will be applied on next system startup",
			"note":    "Currently running on non-Linux system, network configuration changes require manual intervention",
		})
		return
	}

	var req struct {
		InterfaceName string `json:"interface_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	config, err := model.GetNetworkConfigByInterface(req.InterfaceName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Network config not found"})
		return
	}

	go applyLinuxNetworkConfig(config)

	username, _ := c.Get("username")
	model.AddSystemLog("network", "Network config applied for "+req.InterfaceName+" by "+username.(string), "info")

	c.JSON(http.StatusOK, gin.H{"message": "Network configuration applied successfully"})
}

func applyLinuxNetworkConfig(config *model.NetworkConfig) {
	if config.Mode == "dhcp" {
		exec.Command("dhclient", "-r", config.InterfaceName).Run()
		exec.Command("dhclient", config.InterfaceName).Run()
	} else {
		exec.Command("ifconfig", config.InterfaceName, config.IPAddress, "netmask", config.Netmask).Run()
		if config.Gateway != "" {
			exec.Command("route", "add", "default", "gw", config.Gateway).Run()
		}
	}
}

func PingTest(c *gin.Context) {
	var req struct {
		Host string `json:"host" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", req.Host)
	} else {
		cmd = exec.Command("ping", "-c", "4", req.Host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"output":  string(output),
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"output":  string(output),
	})
}

func GetDNSConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"dns_servers": getDNSServers(),
	})
}

func getDNSServers() []string {
	var servers []string
	if runtime.GOOS == "linux" {
		data, err := exec.Command("cat", "/etc/resolv.conf").Output()
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "nameserver") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						servers = append(servers, parts[1])
					}
				}
			}
		}
	}
	return servers
}
