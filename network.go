package main

import (
	"fmt"
)

// NetworkAdapter 网络适配器信息
type NetworkAdapter struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// CurrentNetworkInfo 当前网络信息
type CurrentNetworkInfo struct {
	IP      string
	Gateway string
	DNS     string
	Exists  bool
	Error   string
}

// GetNetworkAdapters 获取所有网络适配器
func GetNetworkAdapters() ([]NetworkAdapter, error) {
	// 使用Windows API,不会有窗口弹出和编码问题
	return GetNetworkAdaptersWinAPI()
}

// GetCurrentConfig 获取当前网络配置
func GetCurrentConfig(interfaceName string) CurrentNetworkInfo {
	// 使用Windows API
	return GetCurrentConfigWinAPI(interfaceName)
}

// GetCurrentConfigName 获取当前配置名称
func GetCurrentConfigName(interfaceName string, configs []NetworkConfig) string {
	current := GetCurrentConfig(interfaceName)

	if !current.Exists {
		return "未检测到配置"
	}

	// 通过网关和IP判断
	for _, cfg := range configs {
		if current.Gateway == cfg.Gateway && current.IP == cfg.IP {
			return cfg.Name
		}
	}

	return fmt.Sprintf("未知配置 (网关: %s)", current.Gateway)
}

// ApplyConfig 应用网络配置
func ApplyConfig(interfaceName string, config NetworkConfig) error {
	// 使用netsh,不会弹窗
	return SetIPConfigWinAPI(interfaceName, config)
}

// IsAdmin 检查是否有管理员权限
func IsAdmin() bool {
	return IsAdminWinAPI()
}
