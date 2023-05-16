package core

import (
	"errors"
	"fmt"
	"github.com/LibCyber/cyber/pkg/clash"
	"github.com/LibCyber/cyber/pkg/util"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrNotAdmin = errors.New("current user is not admin")
)

func EnableTun() error {
	isAdmin, err := util.IsAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return ErrNotAdmin
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %v", err)
	}

	configFilePath := filepath.Join(home, ".cyber", "node", "config.yaml")
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("open config file: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("read config file: %s", err.Error())
	}
	config, err := clash.UnmarshalClashRawConfig(content)
	if err != nil {
		return fmt.Errorf("unmarshal config file: %s", err.Error())
	}

	config.Tun = &clash.TunConfig{
		Enable:              true,
		Stack:               "system",
		AutoRoute:           true,
		AutoDetectInterface: true,
		DnsHijack: []string{
			"1.1.1.1",
			"8.8.8.8",
			"223.5.5.5",
		},
	}

	newContent, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config yaml: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configFilePath, newContent, fileMode)
	if err != nil {
		return fmt.Errorf("writing config file: %v", err)
	}

	_, err = Restart()
	if err != nil {
		return fmt.Errorf("restart cyber-core: %v", err)
	}

	return nil
}

func DisableTun() error {
	isAdmin, err := util.IsAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return ErrNotAdmin
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %v", err)
	}

	configFilePath := filepath.Join(home, ".cyber", "node", "config.yaml")
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("open config file: %s", err.Error())
	}

	content, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("read config file: %s", err.Error())
	}

	config, err := clash.UnmarshalClashRawConfig(content)
	if err != nil {
		return fmt.Errorf("unmarshal config file: %s", err.Error())
	}

	config.Tun = nil

	newContent, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config yaml: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configFilePath, newContent, fileMode)
	if err != nil {
		return fmt.Errorf("writing config file: %v", err)
	}

	_, err = Restart()
	if err != nil {
		return fmt.Errorf("restart cyber-core: %v", err)
	}

	return nil
}

//goland:noinspection GoUnusedExportedFunction
func IsTunEnabled() (bool, error) {
	// 检查配置
	home, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("getting user home dir: %v", err)
	}

	configFilePath := filepath.Join(home, ".cyber", "node", "config.yaml")

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return false, fmt.Errorf("open config file: %s", err.Error())
	}

	content, err := io.ReadAll(configFile)
	if err != nil {
		return false, fmt.Errorf("read config file: %s", err.Error())
	}

	config, err := clash.UnmarshalClashRawConfig(content)
	if err != nil {
		return false, fmt.Errorf("unmarshal config file: %s", err.Error())
	}

	if config.Tun == nil {
		return false, nil
	}

	if !config.Tun.Enable {
		return false, nil
	}

	// 检查虚拟网卡tun是否存在
	_, err = os.Stat("/dev/net/tun")
	if err != nil {
		return false, nil
	}

	return false, nil
}
