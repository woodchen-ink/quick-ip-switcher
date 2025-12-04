package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	SelectedInterface string            `json:"selected_interface"`
	NetworkConfigs    []NetworkConfig   `json:"network_configs"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	SubnetMask string `json:"subnet_mask"`
	Gateway    string `json:"gateway"`
	DNS        string `json:"dns"`
}

var configFile = "config.json"

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(filepath.Dir(exePath), configFile)

	// 如果配置文件不存在,返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			SelectedInterface: "",
			NetworkConfigs:    []NetworkConfig{},
		}, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfig 保存配置
func SaveConfig(config *Config) error {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	configPath := filepath.Join(filepath.Dir(exePath), configFile)

	// 序列化配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	return nil
}

// AddNetworkConfig 添加网络配置
func (c *Config) AddNetworkConfig(config NetworkConfig) {
	c.NetworkConfigs = append(c.NetworkConfigs, config)
}

// RemoveNetworkConfig 删除网络配置
func (c *Config) RemoveNetworkConfig(index int) error {
	if index < 0 || index >= len(c.NetworkConfigs) {
		return fmt.Errorf("索引越界")
	}
	c.NetworkConfigs = append(c.NetworkConfigs[:index], c.NetworkConfigs[index+1:]...)
	return nil
}

// UpdateNetworkConfig 更新网络配置
func (c *Config) UpdateNetworkConfig(index int, config NetworkConfig) error {
	if index < 0 || index >= len(c.NetworkConfigs) {
		return fmt.Errorf("索引越界")
	}
	c.NetworkConfigs[index] = config
	return nil
}

// GetNetworkConfig 获取网络配置
func (c *Config) GetNetworkConfig(index int) (*NetworkConfig, error) {
	if index < 0 || index >= len(c.NetworkConfigs) {
		return nil, fmt.Errorf("索引越界")
	}
	return &c.NetworkConfigs[index], nil
}
