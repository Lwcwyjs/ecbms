package model

import "time"

type NetworkConfig struct {
	ID           int       `json:"id"`
	InterfaceName string   `json:"interface_name"`
	Mode         string    `json:"mode"`
	IPAddress    string    `json:"ip_address"`
	Netmask      string    `json:"netmask"`
	Gateway      string    `json:"gateway"`
	DNSServers   string    `json:"dns_servers"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NetworkInterface struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	HWAddr    string `json:"hwaddr"`
	IPAddress string `json:"ip_address"`
	Netmask   string `json:"netmask"`
}

func GetNetworkConfigByInterface(iface string) (*NetworkConfig, error) {
	config := &NetworkConfig{}
	err := DB.QueryRow(
		`SELECT id, interface_name, mode, ip_address, netmask, gateway, dns_servers, created_at, updated_at 
		 FROM network_configs WHERE interface_name = ? ORDER BY id DESC LIMIT 1`,
		iface,
	).Scan(&config.ID, &config.InterfaceName, &config.Mode, &config.IPAddress, &config.Netmask,
		&config.Gateway, &config.DNSServers, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetAllNetworkConfigs() ([]NetworkConfig, error) {
	rows, err := DB.Query(
		`SELECT id, interface_name, mode, ip_address, netmask, gateway, dns_servers, created_at, updated_at 
		 FROM network_configs ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []NetworkConfig
	for rows.Next() {
		var config NetworkConfig
		err := rows.Scan(&config.ID, &config.InterfaceName, &config.Mode, &config.IPAddress, &config.Netmask,
			&config.Gateway, &config.DNSServers, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func SaveNetworkConfig(iface, mode, ip, netmask, gateway, dns string) error {
	_, err := DB.Exec(
		`INSERT INTO network_configs (interface_name, mode, ip_address, netmask, gateway, dns_servers) 
		 VALUES (?, ?, ?, ?, ?, ?)`,
		iface, mode, ip, netmask, gateway, dns,
	)
	return err
}

func AddSystemLog(logType, message, level string) error {
	_, err := DB.Exec(
		"INSERT INTO system_logs (type, message, level) VALUES (?, ?, ?)",
		logType, message, level,
	)
	return err
}

func GetSystemLogs(limit int) ([]map[string]interface{}, error) {
	rows, err := DB.Query(
		"SELECT id, type, message, level, created_at FROM system_logs ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int
		var logType, message, level string
		var createdAt time.Time
		err := rows.Scan(&id, &logType, &message, &level, &createdAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, map[string]interface{}{
			"id":         id,
			"type":       logType,
			"message":    message,
			"level":      level,
			"created_at": createdAt,
		})
	}
	return logs, nil
}
