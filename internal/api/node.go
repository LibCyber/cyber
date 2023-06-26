package api

import (
	"errors"
	"fmt"
	"github.com/LibCyber/cyber/constant"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/clash"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func (c *Client) DownloadNodes() error {
	info, err := c.GetUserInfo()
	if err != nil {
		// TODO: get from local

		return fmt.Errorf("getting user info: %v", err)
	}

	req, err := http.NewRequest("GET", "https://client-api.libcyber.net/api/client/v1/s/"+info.SubCode, nil)
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberCLI(%s-%s)/%s", runtime.GOOS, runtime.GOARCH, constant.APP_VERSION))
	if err != nil {
		return fmt.Errorf("creating download nodes request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do download nodes request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading download nodes response body: %v", err)
	}

	// 解析 yaml
	newCoreConfig, err := clash.UnmarshalClashRawConfig(body)
	if err != nil {
		return fmt.Errorf("unmarshal config yaml: %v", err)
	}

	// 保存到家目录下的 .cyber/node/config.yaml 中的 accessToken 字段，注意保留文件原来的权限，如果文件不存在，则创建
	usr, err := user.Current()
	if err != nil {
		return errors.New("getting current user")
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	isServiceMode, err := core.IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("checking if service installed: %v", err)
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 检查文件是否存在，不存在则创建并写入，存在则替换 proxies 字段，去重合并rules字段
	_, err = os.Stat(configFilePath)
	if err != nil {
		// 文件不存在，创建
		err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
		if err != nil {
			return fmt.Errorf("creating config file directory: %v", err)
		}
		_, err = os.Create(configFilePath)
		if err != nil {
			return fmt.Errorf("creating config file: %v", err)
		}

		// 写入文件
		err = os.WriteFile(configFilePath, body, 0644)
		if err != nil {
			return fmt.Errorf("writing config file: %v", err)
		}

		return nil
	}

	// 重新获取文件权限
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	// 文件存在，读取文件内容
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("reading config file: %v", err)
	}

	// 解析 yaml
	oldCoreConfig, err := clash.UnmarshalClashRawConfig(fileContent)
	if err != nil {
		return fmt.Errorf("unmarshal config yaml: %v", err)
	}

	// 替换 proxies 字段
	oldCoreConfig.Proxy = newCoreConfig.Proxy

	// 去重合并 rules 字段
	oldCoreConfig.Rule = clash.MergeRules(oldCoreConfig.Rule, newCoreConfig.Rule)

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

func (c *Client) IsNodeExists() (bool, error) {
	usr, err := user.Current()
	if err != nil {
		return false, errors.New("getting current user")
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	isServiceMode, err := core.IsServiceInstalled()
	if err != nil {
		return false, fmt.Errorf("checking if service installed: %v", err)
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 检查文件是否存在，不存在则创建并写入，存在则替换 proxies 字段，去重合并rules字段
	if _, err = os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.New("stat core file: " + err.Error())
	}

	return true, nil
}

func (c *Client) UpdateNodes() error {
	info, err := c.GetUserInfo()
	if err != nil {
		// TODO: get from local

		return fmt.Errorf("getting user info: %v", err)
	}

	req, err := http.NewRequest("GET", "https://client-api.libcyber.net/api/client/v1/s/"+info.SubCode, nil)
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberCLI(%s-%s)/%s", runtime.GOOS, runtime.GOARCH, constant.APP_VERSION))
	if err != nil {
		return fmt.Errorf("creating download nodes request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do download nodes request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading download nodes response body: %v", err)
	}

	// 解析 yaml
	newCoreConfig, err := clash.UnmarshalClashRawConfig(body)
	if err != nil {
		return fmt.Errorf("unmarshal config yaml: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		return errors.New("getting current user")
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	isServiceMode, err := core.IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("checking if service installed: %v", err)
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 获取文件原权限
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	// 打开文件并读取内容
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("reading config file: %v", err)
	}

	// 解析 yaml
	oldCoreConfig, err := clash.UnmarshalClashRawConfig(content)
	if err != nil {
		return fmt.Errorf("unmarshal config yaml: %v", err)
	}

	// 合并配置
	finalCoreConfig := mergeCoreConfig(*oldCoreConfig, *newCoreConfig)

	// 重新编码 yaml
	finalCoreConfigContent, err := yaml.Marshal(finalCoreConfig)
	if err != nil {
		return fmt.Errorf("marshal config yaml: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configFilePath, finalCoreConfigContent, fileMode)
	if err != nil {
		return fmt.Errorf("writing config file: %v", err)
	}

	return nil
}

func mergeCoreConfig(old, new clash.ClashRawConfig) clash.ClashRawConfig {
	// 替换proxies，去重合并rules
	old.Proxy = new.Proxy
	old.Rule = clash.MergeRules(old.Rule, new.Rule)

	return old
}
