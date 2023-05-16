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

type Client struct {
	AccessToken string
}

type Config struct {
	AccessToken string `yaml:"accessToken"`
}

func NewClient() *Client {
	token := extractLocalToken()
	return &Client{
		AccessToken: token,
	}
}

func extractLocalToken() string {
	// 从家目录下的 .cyber/account/config.yaml 中读取 accessToken
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Getting current user:", err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "account", "config.yaml")

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return ""
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return ""
	}

	// base64 解码
	token, err := base64.StdEncoding.DecodeString(config.AccessToken)
	if err != nil {
		return ""
	}
	return string(token)
}

func (c *Client) CheckIn() string {
	return "敬请期待..."
}

type PostMetadata struct {
	OS          string `json:"os"`
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
	FromType    int    `json:"from_type"`
}

func (c *Client) GetUserInfo() (UserProfile, error) {
	if c.AccessToken == "" {
		return UserProfile{}, fmt.Errorf("no access token found, login first using `cyber login`")
	}

	url := "https://client-api.libcyber.net/api/client/v1/profile"
	reqBody := PostMetadata{
		OS:          sysinfo.GetOSInfo(),
		DeviceName:  sysinfo.GetHostname(),
		DeviceModel: sysinfo.GetModel(),
		FromType:    9,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return UserProfile{}, fmt.Errorf("marshalling get user info request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return UserProfile{}, fmt.Errorf("creating get user info request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("LibCyberCLI(%s-%s)/%s", runtime.GOOS, runtime.GOARCH, constant.APP_VERSION))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserProfile{}, fmt.Errorf("do get user info request: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserProfile{}, fmt.Errorf("reading get user info response body: %v", err)
	}

	var userProfile UserProfile
	err = json.Unmarshal(body, &userProfile)
	if err != nil {
		return UserProfile{}, fmt.Errorf("unmarshalling get user info response body: %v", err)
	}

	return userProfile, nil
}
