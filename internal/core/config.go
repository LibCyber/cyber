package core

import (
	"fmt"
	"github.com/LibCyber/cyber/pkg/clash"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
)

func EditNodeConfigWithVi() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %v", err)
	}

	configFilePath := filepath.Join(home, ".cyber", "node", "config.yaml")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("getting service installed: %v", err)
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	cmd := exec.Command("vi", configFilePath)

	// 将Stdin、Stdout和Stderr连接到os.Stdin、os.Stdout和os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("run vi: %s", err.Error())
	}

	return nil
}

func ModifyAPISecret(on bool, secret string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %v", err)
	}

	configFilePath := filepath.Join(home, ".cyber", "node", "config.yaml")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("getting service installed: %v", err)
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 文件存在，读取文件内容
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("reading config file: %v", err)
	}

	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	// 解析 yaml
	oldCoreConfig, err := clash.UnmarshalClashRawConfig(fileContent)
	if err != nil {
		return fmt.Errorf("unmarshal config yaml: %v", err)
	}

	// 替换 secret 字段
	if on {
		oldCoreConfig.Secret = secret
	} else {
		oldCoreConfig.Secret = ""
	}

	// 重新编码 yaml
	newFileContent, err := yaml.Marshal(oldCoreConfig)
	if err != nil {
		return fmt.Errorf("marshal config yaml: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configFilePath, newFileContent, fileMode)
	if err != nil {
		return fmt.Errorf("writing config file: %v", err)
	}

	return nil
}
