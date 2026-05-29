package handler

import (
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"ecbms/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	Hostname        string  `json:"hostname"`
	OS              string  `json:"os"`
	Platform        string  `json:"platform"`
	PlatformVersion string  `json:"platform_version"`
	Uptime          uint64  `json:"uptime"`
	KernelVersion   string  `json:"kernel_version"`
	CPUModel        string  `json:"cpu_model"`
	CPUCores        int     `json:"cpu_cores"`
	TotalMemory     uint64  `json:"total_memory"`
}

type SystemStats struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	UsedMemory   uint64  `json:"used_memory"`
	TotalMemory  uint64  `json:"total_memory"`
	DiskUsage    float64 `json:"disk_usage"`
	UsedDisk     uint64  `json:"used_disk"`
	TotalDisk    uint64  `json:"total_disk"`
	NetworkIn    uint64  `json:"network_in"`
	NetworkOut   uint64  `json:"network_out"`
	Timestamp    int64   `json:"timestamp"`
}

func GetSystemInfo(c *gin.Context) {
	hostInfo, err := host.Info()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get host info"})
		return
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CPU info"})
		return
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get memory info"})
		return
	}

	cpuModel := ""
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	sysInfo := SystemInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		Uptime:          hostInfo.Uptime,
		KernelVersion:   hostInfo.KernelVersion,
		CPUModel:        cpuModel,
		CPUCores:        runtime.NumCPU(),
		TotalMemory:     memInfo.Total,
	}

	c.JSON(http.StatusOK, sysInfo)
}

func GetSystemStats(c *gin.Context) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CPU stats"})
		return
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get memory stats"})
		return
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		diskInfo = &disk.UsageStat{
			UsedPercent: 0,
			Used:        0,
			Total:       1,
		}
	}

	netIO, err := net.IOCounters(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network stats"})
		return
	}

	var netIn, netOut uint64
	if len(netIO) > 0 {
		netIn = netIO[0].BytesRecv
		netOut = netIO[0].BytesSent
	}

	stats := SystemStats{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memInfo.UsedPercent,
		UsedMemory:  memInfo.Used,
		TotalMemory: memInfo.Total,
		DiskUsage:   diskInfo.UsedPercent,
		UsedDisk:    diskInfo.Used,
		TotalDisk:   diskInfo.Total,
		NetworkIn:   netIn,
		NetworkOut:  netOut,
		Timestamp:   time.Now().Unix(),
	}

	c.JSON(http.StatusOK, stats)
}

func GetProcesses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Process list endpoint - requires additional implementation",
	})
}

func RebootSystem(c *gin.Context) {
	if runtime.GOOS == "linux" {
		go func() {
			time.Sleep(2 * time.Second)
			exec.Command("reboot").Run()
		}()
	}

	username, _ := c.Get("username")
	model.AddSystemLog("system", "System reboot initiated by "+username.(string), "warning")

	c.JSON(http.StatusOK, gin.H{"message": "System reboot initiated"})
}

func ShutdownSystem(c *gin.Context) {
	if runtime.GOOS == "linux" {
		go func() {
			time.Sleep(2 * time.Second)
			exec.Command("shutdown", "-h", "now").Run()
		}()
	}

	username, _ := c.Get("username")
	model.AddSystemLog("system", "System shutdown initiated by "+username.(string), "warning")

	c.JSON(http.StatusOK, gin.H{"message": "System shutdown initiated"})
}

func GetSystemLogs(c *gin.Context) {
	limit := 100
	logs, err := model.GetSystemLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
