package app

import (
	"encoding/json"
	"fmt"
	"github.com/LibCyber/cyber/pkg/config"
	"os"
	"path/filepath"
)

func GetCLIConfig() (*config.Config, error) {
	// 检查是否有 ~/.cyber/config.json 目录，没有则创建
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home dir: %s", err.Error())
	}
	cyberConfigPath := filepath.Join(home, ".cyber", "config.json")
	if _, err = os.Stat(cyberConfigPath); os.IsNotExist(err) {
		// 创建文件
		if err = os.MkdirAll(filepath.Dir(cyberConfigPath), 0755); err != nil {
			return nil, fmt.Errorf("create config dir: %s", err.Error())
		}

		// 默认配置
		defaultConfig := config.Config{
			Language: "en",
		}

		// 序列化成json，格式化
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal default config: %s", err.Error())
		}

		// 写入文件
		if err = os.WriteFile(cyberConfigPath, data, 0644); err != nil {
			return nil, fmt.Errorf("write config file: %s", err.Error())
		}
	}

	// 读取文件
	data, err := os.ReadFile(cyberConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %s", err.Error())
	}

	// 反序列化
	var c config.Config
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %s", err.Error())
	}

	return &c, nil
}

func SetLanguage(language string) error {
	// 获取配置
	c, err := GetCLIConfig()
	if err != nil {
		return fmt.Errorf("get config: %s", err.Error())
	}

	// 修改配置
	c.Language = language

	// 序列化成json，格式化
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %s", err.Error())
	}

	// 写入文件
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home dir: %s", err.Error())
	}
	cyberConfigPath := filepath.Join(home, ".cyber", "config.json")
	if err = os.WriteFile(cyberConfigPath, data, 0644); err != nil {
		return fmt.Errorf("write config file: %s", err.Error())
	}

	return nil
}

func Language() string {
	// 获取配置
	c, err := GetCLIConfig()
	if err != nil {
		return "en"
	}

	return c.Language
}
