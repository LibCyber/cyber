package core

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/pkg/country"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotRunning = fmt.Errorf("cyber core is not running")
)

const (
	AUTO_SELECT_PROXY_NAME = "自动选择"
	FALLBACK_PROXY_NAME    = "故障转移"
)

func IsCoreInstalled() (bool, error) {
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return false, errors.New("check service installed: " + err.Error())
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, errors.New("get user home dir: " + err.Error())
	}

	corePath := filepath.Join(homeDir, ".cyber", "core")
	if isServiceMode {
		corePath = filepath.Join("/usr", "bin", "cyber-core")
	}

	if _, err = os.Stat(corePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.New("stat core file: " + err.Error())
	}

	return true, nil
}

func IsNotRunning(err error) bool {
	return errors.Is(err, ErrNotRunning)
}

func Download() error {
	filename := fmt.Sprintf("cyber-core-%s-%s-%s.zip", "v1.0.0", runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(fmt.Sprintf("https://download.libcyber.xyz/clients/cli/%s", filename))
	if err != nil {
		return fmt.Errorf("download core: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// 检查是否有 ~/.cyber/core 目录，没有则创建
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("get current user: %s", err.Error())
	}
	corePath := filepath.Join(usr.HomeDir, ".cyber", "core")
	if _, err = os.Stat(corePath); os.IsNotExist(err) {
		err = os.MkdirAll(corePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create core path: %s", err.Error())
		}
	}

	// 将下载的文件写入到 ~/.cyber/core 目录下
	coreFile, err := os.Create(filepath.Join(corePath, filename))
	if err != nil {
		return fmt.Errorf("create core file: %s", err.Error())
	}

	_, err = io.Copy(coreFile, resp.Body)
	if err != nil {
		return fmt.Errorf("write core file: %s", err.Error())
	}

	// 解压
	err = unzip(filepath.Join(corePath, filename), corePath)
	if err != nil {
		return fmt.Errorf("extract core file: %s", err.Error())
	}

	// 删除压缩包
	err = os.RemoveAll(filepath.Join(corePath, filename))
	if err != nil {
		return fmt.Errorf("remove core zip file: %s", err.Error())
	}

	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("check service installed: %s", err.Error())
	}

	if isServiceMode {
		// copy ~/.cyber/core/cyber-core to /usr/bin/cyber-core
		err = util.CopyFile(corePath, filepath.Join("/usr", "bin", "cyber-core"))
		if err != nil {
			return err
		}
	}

	return nil
}

func unzip(src, dest string) error {
	// 打开zip文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer r.Close()

	// 遍历zip文件中的每个文件/目录
	for _, f := range r.File {
		// 计算解压后的文件路径
		fpath := filepath.Join(dest, f.Name)

		// 创建目录
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// 创建解压后的文件
		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			return err
		}

		// 解压文件
		err = extractFile(f, fpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(f *zip.File, destPath string) error {
	// 打开zip文件中的文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer rc.Close()

	// 创建目标文件
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer outFile.Close()

	// 将文件内容复制到目标文件
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	return nil
}

func GetProxyPort() (int, int, error) {
	usr, err := user.Current()
	if err != nil {
		return 0, 0, fmt.Errorf("get current user: %s", err.Error())
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return 0, 0, fmt.Errorf("check service installed: %s", err.Error())
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 获取配置文件中的端口号
	httpPort, socksPort, _, err := getPortInConfig(filepath.Join(configFilePath))
	if err != nil {
		return 0, 0, err
	}

	return httpPort, socksPort, nil
}

func getPortInConfig(configPath string) (int, int, int, error) {
	// 打开配置文件，解析yaml
	configFile, err := os.Open(configPath)
	if err != nil {
		return 0, 0, 0, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	config := make(map[any]any)
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return 0, 0, 0, err
	}

	// 获取端口号
	port, ok := config["port"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("port not found in config file")
	}

	socksPort, ok := config["socks-port"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("socks-port not found in config file")
	}

	externalController, ok := config["external-controller"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("external-controller not found in config file")
	}
	// 解析出端口号
	_, apiPortStr, err := net.SplitHostPort(externalController.(string))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse external-controller: %s", err.Error())
	}

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("convert api port: %s", err.Error())
	}

	return port.(int), socksPort.(int), apiPort, nil
}

func setPortInConfig(configPath string, port, socksPort int, apiPort int) error {
	// 打开配置文件，解析yaml
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	config := make(map[any]any)
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return err
	}

	// 修改端口号
	config["port"] = port
	config["socks-port"] = socksPort

	addr, _, err := net.SplitHostPort(config["external-controller"].(string))
	if err != nil {
		return fmt.Errorf("parse external-controller: %s", err.Error())
	}

	config["external-controller"] = net.JoinHostPort(addr, strconv.Itoa(apiPort))

	// 将修改后的配置写入到配置文件
	configFile, err = os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	err = yaml.NewEncoder(configFile).Encode(&config)
	if err != nil {
		return err
	}

	return nil
}

func Start() (int, error) {
	exist, pid, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return 0, fmt.Errorf("check user process: %s", err.Error())
	}

	if exist {
		return 0, fmt.Errorf("cyber core is already running, pid: %d", pid)
	}

	usr, err := user.Current()
	if err != nil {
		return 0, fmt.Errorf("get current user: %s", err.Error())
	}

	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return 0, fmt.Errorf("check service installed: %s", err.Error())
	}

	corePath := filepath.Join(usr.HomeDir, ".cyber", "core")
	if isServiceMode {
		corePath = filepath.Join("/usr", "bin", "cyber-core")
	}
	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// check if cyber-core directory exists
	if isServiceMode {
		if _, err = os.Stat(corePath); os.IsNotExist(err) {
			return 0, fmt.Errorf("cyber core is not installed, please install it first using `cyber core download` and install service again using `cyber core enable`")
		}
	} else {
		if _, err = os.Stat(corePath); os.IsNotExist(err) {
			return 0, fmt.Errorf("cyber core is not installed, please install it first using `cyber core download`")
		}
	}

	// 获取配置文件中的端口号
	httpPort, socksPort, apiPort, err := getPortInConfig(filepath.Join(configFilePath))
	if err != nil {
		return 0, fmt.Errorf("get port: %s", err.Error())
	}

	var changed bool
	// 检查端口是否被占用
	if !checkPortAvailable(httpPort) {
		httpPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if httpPort < 0 {
			return 0, fmt.Errorf("no available port for http proxy")
		}

		changed = true
	}

	if !checkPortAvailable(socksPort) {
		socksPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if socksPort < 0 {
			return 0, fmt.Errorf("no available port for socks proxy")
		}

		changed = true
	}

	if !checkPortAvailable(apiPort) {
		apiPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if apiPort < 0 {
			return 0, fmt.Errorf("no available port for api")
		}

		changed = true
	}

	if changed {
		// 修改配置文件中的端口号
		err = setPortInConfig(configFilePath, httpPort, socksPort, apiPort)
	}

	if isServiceMode {
		// execute `systemctl start cyber-core`
		cmd := exec.Command("systemctl", "start", "cyber-core")
		err := cmd.Run()
		if err != nil {
			return 0, fmt.Errorf("start cyber core service: %s", err.Error())
		}

		// get pid
		cmd = exec.Command("systemctl", "show", "--property=MainPID", "cyber-core")
		output, err := cmd.Output()
		if err != nil {
			return 0, fmt.Errorf("get cyber core service pid: %s", err.Error())
		}

		// parse pid
		re := regexp.MustCompile(`MainPID=(\d+)`)
		match := re.FindSubmatch(output)
		if match == nil {
			return 0, fmt.Errorf("parse cyber core service pid: %s", err.Error())
		}

		mainPid, err := strconv.Atoi(string(match[1]))
		if err != nil {
			return 0, fmt.Errorf("convert cyber core service pid: %s", err.Error())
		}

		pid = mainPid
	} else {

		cmd := exec.Command(filepath.Join(corePath, "cyber-core"), "-d", corePath, "-f", configFilePath)

		// 创建日志文件，以覆写的方式打开
		logFile, err := os.OpenFile(filepath.Join(corePath, "coreOut.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to open log file: %v", err)
		}
		//goland:noinspection GoUnhandledErrorResult
		defer logFile.Close()

		errLogFile, err := os.OpenFile(filepath.Join(corePath, "coreErr.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to open error log file: %v", err)
		}
		//goland:noinspection GoUnhandledErrorResult
		defer errLogFile.Close()

		// 重定向程序输出到日志文件
		cmd.Stdout = logFile
		cmd.Stderr = errLogFile

		// 将进程分离到新的进程组，以避免终端信号影响子进程
		cmd.SysProcAttr = procAttrWithNewProcessGroup()

		// 启动子进程
		err = cmd.Start()
		if err != nil {
			return 0, fmt.Errorf("failed to start core: %v", err)
		}
		pid = cmd.Process.Pid
	}

	timeout := time.Second * 10
	startTime := time.Now()
	for {
		if time.Now().Sub(startTime) > timeout {
			return 0, fmt.Errorf("timeout")
		}

		err = pingCore()
		if err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	return pid, nil
}

//goland:noinspection GoUnusedFunction
func killProcessesByName(name string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/IM", name, "/F", "/T")
	case "linux", "darwin", "freebsd":
		// 获取当前用户的用户ID
		currentUser, err := user.Current()
		if err != nil {
			return err
		}
		uid, err := strconv.Atoi(currentUser.Uid)
		if err != nil {
			return err
		}

		// 查找进程
		pgrepCmd := exec.Command("pgrep", "-f", "-u", strconv.Itoa(uid), name)
		output, err := pgrepCmd.Output()
		if err != nil {
			return err
		}

		// 解析进程ID
		pids := strings.Split(strings.TrimSpace(string(output)), "\n")

		// 杀死进程
		for _, pid := range pids {
			cmd = exec.Command("kill", pid)
			err = cmd.Run()
			if err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func killProcessesByPid(pid int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F", "/T")
	case "linux", "darwin", "freebsd":
		cmd = exec.Command("kill", strconv.Itoa(pid))
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func Stop() error {
	exist, pid, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return nil
	}

	//err = killProcessesByName(processName)
	//if err != nil {
	//	return err
	//}

	// close all connections first
	err = closeAllConnections()
	if err != nil {
		return fmt.Errorf("close all connections: %s", err.Error())
	}

	time.Sleep(time.Second) // wait for connections to close

	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("check service installed: %s", err.Error())
	}

	if isServiceMode {
		cmd := exec.Command("systemctl", "stop", "cyber-core")
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("stop cyber core service: %s", err.Error())
		}
	} else {
		err = killProcessesByPid(pid)
		if err != nil {
			return fmt.Errorf("kill process: %s", err.Error())
		}
	}

	timeout := time.Second * 10
	startTime := time.Now()
	for {
		exist, _, _, err = checkUserProcess("cyber-core")
		if err != nil {
			return fmt.Errorf("check user process: %s", err.Error())
		}

		if !exist {
			break
		}

		if time.Now().Sub(startTime) > timeout {
			return fmt.Errorf("timeout")
		}

		time.Sleep(time.Second)
	}

	return nil
}

func Restart() (int, error) {
	err := Stop()
	if err != nil && !IsNotRunning(err) {
		return 0, fmt.Errorf("stop cyber core: %s", err.Error())
	}

	pid, err := Start()
	if err != nil {
		return 0, fmt.Errorf("start cyber core: %s", err.Error())
	}
	return pid, nil
}

// Status returns the pid of cyber core process if it is running, otherwise returns 0.
func Status() (int, error) {
	exist, pid, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return 0, fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return 0, nil
	}

	return pid, nil
}

//func CheckUpdate() {
//
//}

func checkPortAvailable(port int) bool {
	if port < 1 || port > 65534 {
		return false
	}
	addr := ":"
	l, err := net.Listen("tcp", addr+strconv.Itoa(port))
	if err != nil {
		//log.Warnln("check port fail 0.0.0.0:%d", port)
		return false
	}
	_ = l.Close()

	addr = "127.0.0.1:"
	l, err = net.Listen("tcp", addr+strconv.Itoa(port))
	if err != nil {
		//log.Warnln("check port fail 127.0.0.1:%d", port)
		return false
	}
	_ = l.Close()
	//log.Infoln("check port %d success", port)
	return true
}

// 获取可用端口，从10000-65534，不包括maskPort，如果没有可用端口则返回-1
func getAvailablePort(maskPorts []int) int {
	for i := 10000; i <= 65534; i++ {
		if !checkPortAvailable(i) {
			continue
		}
		if maskPorts != nil {
			isMask := false
			for _, v := range maskPorts {
				if v == i {
					isMask = true
					break
				}
			}
			if isMask {
				continue
			}
		}
		return i
	}
	return -1
}

type ConfigRestfulResponse struct {
	Port           int           `json:"port"`
	SocksPort      int           `json:"socks-port"`
	RedirPort      int           `json:"redir-port"`
	TproxyPort     int           `json:"tproxy-port"`
	MixedPort      int           `json:"mixed-port"`
	Authentication []interface{} `json:"authentication"`
	AllowLan       bool          `json:"allow-lan"`
	BindAddress    string        `json:"bind-address"`
	Mode           string        `json:"mode"`
	LogLevel       string        `json:"log-level"`
	Ipv6           bool          `json:"ipv6"`
}

func GetConfigs() (*ConfigRestfulResponse, error) {
	externalController, err := getExternalController()
	if err != nil {
		return nil, fmt.Errorf("get external controller: %s", err.Error())
	}
	resp, err := http.Get(fmt.Sprintf("http://%s/configs", externalController))
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return nil, fmt.Errorf("get configs: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %s", err.Error())
	}

	var configs ConfigRestfulResponse
	err = yaml.Unmarshal(body, &configs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body: %s", err.Error())
	}

	return &configs, nil
}

func getSelector() (string, error) {
	c, err := GetConfigs()
	if err != nil {
		return "", err
	}

	if c.Mode == "global" {
		return "GLOBAL", nil
	}
	if c.Mode == "rule" {
		return "SELECT", nil
	}

	return "", fmt.Errorf("unknown mode: %s", c.Mode)
}

func ChangeNode(nodeName string) error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	nodes, _, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes: %s", err.Error())
	}

	// 检查节点是否存在
	nodeExist := false
	for _, v := range nodes {
		if v.Name == nodeName {
			nodeExist = true
			break
		}
	}
	if !nodeExist {
		return fmt.Errorf("node %s not exist", nodeName)
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller: %s", err.Error())
	}

	selector, err := getSelector()
	if err != nil {
		return fmt.Errorf("get selector: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/proxies/%s", externalController, selector), nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return fmt.Errorf("new request: %s", err.Error())
	}

	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, nodeName)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("change node: %s", resp.Status)
	}
	return nil
}

func ChangeMode(mode string) error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller: %s", err.Error())
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("http://%s/configs", externalController), nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return fmt.Errorf("new request: %s", err.Error())
	}

	if mode != "global" && mode != "rule" {
		return fmt.Errorf("unknown mode: %s", mode)
	}
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"mode": "%s"}`, mode)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("change mode: %s", resp.Status)
	}

	return nil
}

type DelayTestResponse struct {
	Delay   int    `json:"delay"`
	Message string `json:"message"`
}

func BenchmarkNode() error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller: %s", err.Error())
	}

	nodes, _, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes: %s", err.Error())
	}

	for _, node := range nodes {
		var delayTestResponse *DelayTestResponse
		delayTestResponse, statusCode, status, err := sendDelayTestRequest(externalController, node)
		if err != nil {
			return fmt.Errorf("send delay test request: %s", err.Error())
		}

		if statusCode != http.StatusOK {
			if app.Language() == "zh" {
				fmt.Printf("*** 正在测试节点 %s, 状态: %s, 信息: %s\n", node.Name, status, delayTestResponse.Message)
			} else {
				fmt.Printf("*** Benchmarking node %s, status: %s, message: %s\n", node.Name, status, delayTestResponse.Message)
			}
			continue
		}

		if app.Language() == "zh" {
			fmt.Printf("正在测试节点 %s, 状态: %s, 延迟: %d毫秒\n", node.Name, status, delayTestResponse.Delay)
		} else {
			fmt.Printf("Benchmarking node %s, status: %s, delay: %dms\n", node.Name, status, delayTestResponse.Delay)
		}
	}

	return nil
}

func closeAllConnections() error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller: %s", err.Error())
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/connections", externalController), nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return fmt.Errorf("new request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %s", err.Error())
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("close all connections: %s", resp.Status)
	}

	return nil
}

func pingCore() error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller: %s", err.Error())
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s", externalController), nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return fmt.Errorf("new request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %s", err.Error())
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping core: %s", resp.Status)
	}

	return nil
}

func sendDelayTestRequest(externalController string, node Node) (*DelayTestResponse, int, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/proxies/%s/delay?timeout=%d&url=%s", externalController, node.Name, 3000, "http://www.gstatic.com/generate_204"), nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, 0, "", fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return nil, 0, "", fmt.Errorf("new request: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, "", fmt.Errorf("do request: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", fmt.Errorf("read response body: %s", err.Error())
	}

	var delayTestResponse DelayTestResponse
	err = json.Unmarshal(body, &delayTestResponse)
	if err != nil {
		return nil, 0, "", fmt.Errorf("unmarshal response body: %s", err.Error())
	}

	return &delayTestResponse, resp.StatusCode, resp.Status, nil
}

type Node struct {
	Name  string
	Type  string // ss, ssr, v2ray, trojan
	Delay int
}

type ProxySelectorResponse struct {
	All     []string `json:"all"`
	History []any    `json:"history"`
	Name    string   `json:"name"`
	Now     string   `json:"now"`
	Type    string   `json:"type"`
	Udp     bool     `json:"udp"`
}

type ProxyRestfulResponse struct {
	History []struct {
		Time  time.Time `json:"time"`
		Delay int       `json:"delay"`
	} `json:"history"`
	Name string `json:"name"`
	Type string `json:"type"`
	Udp  bool   `json:"udp"`
}

func ListNodes() error {
	exist, _, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process: %s", err.Error())
	}

	if !exist {
		return ErrNotRunning
	}

	nodes, currentNode, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes: %s", err.Error())
	}

	nodeOfAutoSelect, err := getAutoSelectNodeName()
	if err != nil {
		return fmt.Errorf("get auto select node name: %s", err.Error())
	}

	nodeOfFallback, err := getFallbackNodeName()
	if err != nil {
		return fmt.Errorf("get fallback node name: %s", err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	if app.Language() == "zh" {
		table.SetHeader([]string{"当前选择", "节点名称", "节点类型", "延迟", "位置"})
	} else {
		table.SetHeader([]string{"Current Selected", "Node Name", "Node Type", "Delay", "Location"})
	}

	for _, node := range nodes {
		var data []string
		if node.Name == currentNode {
			data = append(data, "*")
		} else {
			data = append(data, "")
		}
		nodeName := node.Name
		if node.Name == AUTO_SELECT_PROXY_NAME {
			nodeName = fmt.Sprintf("%s (%s)", node.Name, nodeOfAutoSelect)
		} else if node.Name == FALLBACK_PROXY_NAME {
			nodeName = fmt.Sprintf("%s (%s)", node.Name, nodeOfFallback)
		}
		if node.Delay == 0 {
			data = append(data, nodeName, node.Type, "N/A", getCountryCode(nodeName))
			table.Append(data)
			continue
		} else {
			if app.Language() == "zh" {
				data = append(data, nodeName, node.Type, fmt.Sprintf("%d毫秒", node.Delay), getCountryCode(nodeName))
			} else {
				data = append(data, nodeName, node.Type, fmt.Sprintf("%dms", node.Delay), getCountryCode(nodeName))
			}
		}
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Render() // Send output

	return nil
}

func getNodes() ([]Node, string, error) {
	externalController, err := getExternalController()
	if err != nil {
		return nil, "", fmt.Errorf("get external controller: %s", err.Error())
	}

	selector, err := getSelector()
	if err != nil {
		return nil, "", fmt.Errorf("get selector: %s", err.Error())
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, selector))
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, "", fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return nil, "", fmt.Errorf("get proxies: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %s", err.Error())
	}

	var proxies ProxySelectorResponse

	err = json.Unmarshal(body, &proxies)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal body: %s", err.Error())
	}

	var nodes []Node
	for _, nodeName := range proxies.All {
		var node *ProxyRestfulResponse
		node, err = getProxy(externalController, nodeName)
		if err != nil {
			return nil, "", fmt.Errorf("get proxy: %s", err.Error())
		}

		var delay int
		if len(node.History) == 0 {
			delay = 0
		} else {
			delay = node.History[0].Delay
		}
		nodes = append(nodes, Node{
			Name:  node.Name,
			Type:  node.Type,
			Delay: delay,
		})
	}

	return nodes, proxies.Now, nil
}

func getAutoSelectNodeName() (string, error) {
	externalController, err := getExternalController()
	if err != nil {
		return "", fmt.Errorf("get external controller: %s", err.Error())
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, AUTO_SELECT_PROXY_NAME))
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return "", fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return "", fmt.Errorf("get proxies: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %s", err.Error())
	}

	var proxy ProxySelectorResponse

	err = json.Unmarshal(body, &proxy)
	if err != nil {
		return "", fmt.Errorf("unmarshal body: %s", err.Error())
	}

	return proxy.Now, nil
}

func getFallbackNodeName() (string, error) {
	externalController, err := getExternalController()
	if err != nil {
		return "", fmt.Errorf("get external controller: %s", err.Error())
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, FALLBACK_PROXY_NAME))
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return "", fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return "", fmt.Errorf("get proxies: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %s", err.Error())
	}

	var proxy ProxySelectorResponse

	err = json.Unmarshal(body, &proxy)
	if err != nil {
		return "", fmt.Errorf("unmarshal body: %s", err.Error())
	}

	return proxy.Now, nil
}

func getProxy(externalController string, nodeName string) (*ProxyRestfulResponse, error) {
	nodeResp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, nodeName))
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, fmt.Errorf("cyber core api endpoint is not available, please check if cyber core is running")
		}
		return nil, fmt.Errorf("get proxy: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer nodeResp.Body.Close()

	nodeBody, err := io.ReadAll(nodeResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read node body: %s", err.Error())
	}

	var node ProxyRestfulResponse
	err = json.Unmarshal(nodeBody, &node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal node body: %s", err.Error())
	}

	return &node, nil
}

func getExternalController() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("get current user: %s", err.Error())
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return "", fmt.Errorf("check if service is installed: %s", err.Error())
	}
	if isServiceMode {
		configFilePath = filepath.Join("/etc", "cyber-core", "config.yaml")
	}

	// 打开配置文件，解析yaml
	configFile, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("config file not found, login first using `cyber login` and download nodes using `cyber node download`")
		}
		return "", fmt.Errorf("open config file: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	config := make(map[any]any)
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return "", fmt.Errorf("decode config file: %s", err.Error())
	}

	externalController, ok := config["external-controller"]
	if !ok {
		return "", fmt.Errorf("external-controller not found")
	}

	return externalController.(string), nil
}

//goland:noinspection GoUnusedParameter
func getCountryCode(nodeName string) string {
	countryCode, err := country.ParseCountry(nodeName)
	if err != nil {
		return "N/A"
	}

	return countryCode
}

func ShowLogFollow() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("get current user: %s", err.Error())
	}
	coreOutLogPath := filepath.Join(usr.HomeDir, ".cyber", "core", "coreOut.log")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("check service installed: %s", err.Error())
	}

	if isServiceMode {
		// journalctl -u cyber-core -f
		cmd := exec.Command("journalctl", "-u", "cyber-core", "-f")
		// 获取输出对象，可以从该对象中读取输出结果
		stdout, _ := cmd.StdoutPipe()
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("start journalctl: %s", err.Error())
		}

		// 创建一个扫描器并扫描stdout
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)

		// 循环打印扫描器扫描到的内容，直到扫描器停止
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}

		err = cmd.Wait()
		if err != nil {
			return errors.New("journalctl exited unexpectedly")
		}

		return nil
	}

	file, err := os.Open(coreOutLogPath)
	if err != nil {
		return fmt.Errorf("open coreOut.log: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	reader := io.NewSectionReader(file, 0, 1<<63-1)
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read coreOut.log: %s", err.Error())
		}

		if n > 0 {
			fmt.Print(string(buf[:n]))
		}

		if err == io.EOF {
			time.Sleep(time.Second) // 等待一秒后再次读取，以检测是否有新的数据
			_, err := reader.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("seek coreOut.log: %s", err.Error())
			} // 重置读取位置，继续读取新的数据
		}
	}
}

func ShowLog() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("get current user: %s", err.Error())
	}
	coreOutLogPath := filepath.Join(usr.HomeDir, ".cyber", "core", "coreOut.log")
	isServiceMode, err := IsServiceInstalled()
	if err != nil {
		return fmt.Errorf("check service installed: %s", err.Error())
	}

	if isServiceMode {
		// journalctl -u cyber-core
		cmd := exec.Command("journalctl", "-u", "cyber-core")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("start journalctl: %s", err.Error())
		}

		err = cmd.Wait()
		if err != nil {
			return errors.New("journalctl exited unexpectedly")
		}

		return nil
	}

	file, err := os.Open(coreOutLogPath)
	if err != nil {
		return fmt.Errorf("open coreOut.log: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read coreOut.log: %s", err.Error())
		}

		if n > 0 {
			fmt.Print(string(buf[:n]))
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func EnableService() error {
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS != "linux" {
		return ErrNotLinux
	}

	// execute systemctl enable cyber-core and print output
	cmd := exec.Command("systemctl", "enable", "cyber-core")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}

func DisableService() error {
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS != "linux" {
		return ErrNotLinux
	}

	// execute systemctl disable cyber-core
	cmd := exec.Command("systemctl", "disable", "cyber-core")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Print(string(output))

	return nil
}
