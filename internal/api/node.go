package api

import (
	"fmt"
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
		// TODO get from local
		return fmt.Errorf("error getting user info: %v", err)
	}

	req, err := http.NewRequest("GET", "https://client-api.libcyber.net/api/client/v1/s/"+info.SubCode, nil)
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberDesktop(%s-%s)/1.0.0", runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return fmt.Errorf("error creating download nodes request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error do download nodes request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading download nodes response body: %v", err)
	}

	// 解析 yaml
	newCoreConfig, err := clash.UnmarshalClashRawConfig(body)
	if err != nil {
		return fmt.Errorf("error unmarshal config yaml: %v", err)
	}

	// 保存到家目录下的 .cyber/node/config.yaml 中的 accessToken 字段，注意保留文件原来的权限，如果文件不存在，则创建
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	// 检查文件是否存在，不存在则创建并写入，存在则替换 proxies 字段，去重合并rules字段
	_, err = os.Stat(configFilePath)
	if err != nil {
		// 文件不存在，创建
		err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
		if err != nil {
			return fmt.Errorf("error creating config file directory: %v", err)
		}
		_, err = os.Create(configFilePath)
		if err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}

		// 写入文件
		err = os.WriteFile(configFilePath, body, 0644)
		if err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		return nil
	}

	// 重新获取文件权限
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("error getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	// 文件存在，读取文件内容
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// 解析 yaml
	oldCoreConfig, err := clash.UnmarshalClashRawConfig(fileContent)
	if err != nil {
		return fmt.Errorf("error unmarshal config yaml: %v", err)
	}

	// 替换 proxies 字段
	oldCoreConfig.Proxy = newCoreConfig.Proxy

	// 去重合并 rules 字段
	oldCoreConfig.Rule = clash.MergeRules(oldCoreConfig.Rule, newCoreConfig.Rule)

	// 重新编码 yaml
	newFileContent, err := yaml.Marshal(oldCoreConfig)
	if err != nil {
		return fmt.Errorf("error marshal config yaml: %v", err)
	}

	// 写入文件
	err = os.WriteFile(configFilePath, newFileContent, fileMode)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
