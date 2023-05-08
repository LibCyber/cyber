package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LibCyber/cyber/constant"
	"github.com/LibCyber/cyber/pkg/sysinfo"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	OS          string `json:"os"`
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
	FromType    int    `json:"from_type"`
}

type UserProfile struct {
	Id               int      `json:"id"`
	Nickname         string   `json:"nickname"`
	Account          string   `json:"account"`
	Port             int      `json:"port"`
	Passwd           string   `json:"passwd"`
	Uuid             string   `json:"uuid"`
	TransferEnable   int      `json:"transfer_enable"`
	U                int      `json:"u"`
	D                int      `json:"d"`
	T                int      `json:"t"`
	Enable           int      `json:"enable"`
	SpeedLimit       int      `json:"speed_limit"`
	Credit           int      `json:"credit"`
	Balance          int      `json:"balance"`
	ExpiredAt        string   `json:"expired_at"`
	BanTime          int      `json:"ban_time"`
	Level            int      `json:"level"`
	IsAdmin          int      `json:"is_admin"`
	Group            any      `json:"group"`
	LastLogin        int      `json:"last_login"`
	ResetTime        int      `json:"reset_time"`
	InviteNum        any      `json:"invite_num"`
	Status           int      `json:"status"`
	Sublink          string   `json:"sublink"`
	SubCode          string   `json:"sub_code"`
	Clashsub         string   `json:"clashsub"`
	ClashSublinkPath string   `json:"clash_sublink_path"`
	LabelList        []string `json:"label_list"`
	TopLabel         string   `json:"top_label"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	Code        int         `json:"code"`
	ExpiredAt   int         `json:"expired_at"`
	ExpiresIn   int         `json:"expires_in"`
	Id          int         `json:"id"`
	Message     string      `json:"message"`
	Ret         int         `json:"ret"`
	Success     bool        `json:"success"`
	TokenType   string      `json:"token_type"`
	User        UserProfile `json:"user"`
	Username    string      `json:"username"`
}

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (c *Client) Login(username string, password string) (string, error) {
	url := "https://client-api.libcyber.net/api/client/v1/login"
	loginRequest := LoginRequest{
		Username:    username,
		Password:    password,
		OS:          sysinfo.GetOSInfo(),
		DeviceName:  sysinfo.GetHostname(),
		DeviceModel: sysinfo.GetModel(),
		FromType:    9,
	}

	jsonData, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("error marshalling login request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating login request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberCLI(%s-%s)/%s", runtime.GOOS, runtime.GOARCH, constant.APP_VERSION))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error do login request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading login response: %v", err)
	}

	//fmt.Println("login response: ", string(body))
	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("error decoding login response: %v", err)
		}

		return "", fmt.Errorf("login failed: %v", errorResponse.ErrorMessage)
	}

	if !loginResponse.Success {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", fmt.Errorf("error decoding login response: %v, %v", err, loginResponse.Message)
		}

		return "", fmt.Errorf("login failed: %v", errorResponse.ErrorMessage)
	}

	err = c.saveCredential(loginResponse.AccessToken)
	if err != nil {
		return "", fmt.Errorf("error saving credential: %v", err)
	}

	return loginResponse.AccessToken, nil
}

func (c *Client) saveCredential(token string) error {
	// 保存到家目录下的 .cyber/account/config.yaml 中的 accessToken 字段，注意保留文件原来的权限，如果文件不存在，则创建
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "account", "config.yaml")
	// 检查文件是否存在
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
	}
	// 重新获取文件权限
	fileInfo, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("error getting config file info: %v", err)
	}
	fileMode := fileInfo.Mode()

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}

	config.AccessToken = base64.StdEncoding.EncodeToString([]byte(token))

	configData, err = yaml.Marshal(config)

	// base64编码后再写入文件
	err = os.WriteFile(configFilePath, []byte(configData), fileMode)

	return nil
}

func (c *Client) removeCredential() error {
	// 删除家目录下的 .cyber/ 目录
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber")
	// 检查文件是否存在
	_, err = os.Stat(configFilePath)
	if err != nil {
		return nil
	}

	err = os.RemoveAll(configFilePath)
	if err != nil {
		return fmt.Errorf("error removing config file directory: %v", err)
	}

	return nil
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) Logout() error {
	url := "https://client-api.libcyber.net/api/client/v1/logout"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating logout request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberCLI(%s-%s)/%s", runtime.GOOS, runtime.GOARCH, constant.APP_VERSION))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error do logout request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading logout response: %v", err)
	}

	var logoutResponse LogoutResponse
	err = json.Unmarshal(body, &logoutResponse)
	if err != nil {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding logout response: %v", err)
		}

		return fmt.Errorf("logout failed: %v", errorResponse.ErrorMessage)
	}

	if !logoutResponse.Success {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding logout response: %v, %v", err, logoutResponse.Message)
		}

		return fmt.Errorf("logout failed: %v", errorResponse.ErrorMessage)
	}

	err = c.removeCredential()
	if err != nil {
		return fmt.Errorf("error removing credential: %v", err)
	}

	return nil
}
