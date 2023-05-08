package core

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/process"
	"gopkg.in/yaml.v3"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Download() error {
	// 从 download.libcyber.xyz/clients/cli/V1.0.0/cyber-v[version]-[GOOS]-[GOARCH].zip 下载
	filename := fmt.Sprintf("cyber-core-%s-%s-%s.zip", "v1.0.0", runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(fmt.Sprintf("https://download.libcyber.xyz/clients/cli/%s", filename))
	if err != nil {
		return fmt.Errorf("download core error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// 检查是否有 ~/.cyber/core 目录，没有则创建
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("get current user error: %s", err.Error())
	}
	corePath := filepath.Join(usr.HomeDir, ".cyber", "core")
	if _, err = os.Stat(corePath); os.IsNotExist(err) {
		err = os.MkdirAll(corePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create core path error: %s", err.Error())
		}
	}

	// 将下载的文件写入到 ~/.cyber/core 目录下
	coreFile, err := os.Create(filepath.Join(corePath, filename))
	if err != nil {
		return fmt.Errorf("create core file error: %s", err.Error())
	}

	_, err = io.Copy(coreFile, resp.Body)
	if err != nil {
		return fmt.Errorf("write core file error: %s", err.Error())
	}

	// 解压
	err = unzip(filepath.Join(corePath, filename), corePath)
	if err != nil {
		return fmt.Errorf("extract core file error: %s", err.Error())
	}

	// 删除压缩包
	err = os.Remove(filepath.Join(corePath, filename))
	if err != nil {
		return fmt.Errorf("remove core zip file error: %s", err.Error())
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
		return 0, 0, fmt.Errorf("get current user error: %s", err.Error())
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")

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
		return 0, 0, 0, fmt.Errorf("parse external-controller error: %s", err.Error())
	}

	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("convert api port error: %s", err.Error())
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
		return fmt.Errorf("parse external-controller error: %s", err.Error())
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

func Start() error {
	exist, pid, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if exist {
		return fmt.Errorf("cyber-core is already running, pid: %d", pid)
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("get current user error: %s", err.Error())
	}

	corePath := filepath.Join(usr.HomeDir, ".cyber", "core")
	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")

	// 获取配置文件中的端口号
	httpPort, socksPort, apiPort, err := getPortInConfig(filepath.Join(configFilePath))
	if err != nil {
		return fmt.Errorf("get port error: %s", err.Error())
	}

	var changed bool
	// 检查端口是否被占用
	if !checkPortAvailable(httpPort) {
		httpPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if httpPort < 0 {
			return fmt.Errorf("no available port for http proxy")
		}

		changed = true
	}

	if !checkPortAvailable(socksPort) {
		socksPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if socksPort < 0 {
			return fmt.Errorf("no available port for socks proxy")
		}

		changed = true
	}

	if !checkPortAvailable(apiPort) {
		apiPort = getAvailablePort([]int{httpPort, socksPort, apiPort})
		if apiPort < 0 {
			return fmt.Errorf("no available port for api")
		}

		changed = true
	}

	if changed {
		// 修改配置文件中的端口号
		err = setPortInConfig(configFilePath, httpPort, socksPort, apiPort)
	}

	cmd := exec.Command(filepath.Join(corePath, "cyber-core"), "-d", corePath, "-f", configFilePath)

	// 创建日志文件，以覆写的方式打开
	logFile, err := os.OpenFile(filepath.Join(corePath, "coreOut.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer logFile.Close()

	errLogFile, err := os.OpenFile(filepath.Join(corePath, "coreErr.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer errLogFile.Close()

	// 重定向程序输出到日志文件
	cmd.Stdout = logFile
	cmd.Stderr = errLogFile

	// 将进程分离到新的进程组，以避免终端信号影响子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// 启动子进程
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start core: %v", err)
	}

	fmt.Printf("Started background process with PID %d\n", cmd.Process.Pid)
	return nil
}

func killProcessesByName(name string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/IM", "*"+name+"*", "/F", "/T")
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

func Stop() error {
	exist, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	err = killProcessesByName("cyber-core")
	if err != nil {
		return err
	}

	return nil
}

func Status() error {
	exist, pid, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	fmt.Println("cyber-core is running, pid:", pid)
	return nil
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

func checkUserProcess(processName string) (bool, int, error) {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		return false, 0, err
	}

	// 将当前用户的UID转换为整数
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false, 0, err
	}

	// 获取系统中所有运行中的进程
	processes, err := process.Processes()
	if err != nil {
		return false, 0, err
	}

	// 遍历所有进程，检查进程名和用户ID
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if name == processName {
			uids, err := p.Uids()
			if err != nil || len(uids) == 0 {
				continue
			}

			if uids[0] == int32(uid) {
				return true, int(p.Pid), nil
			}
		}
	}

	return false, 0, nil
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
		return nil, fmt.Errorf("get external controller error: %s", err.Error())
	}
	resp, err := http.Get(fmt.Sprintf("http://%s/configs", externalController))
	if err != nil {
		return nil, fmt.Errorf("get configs error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %s", err.Error())
	}

	var configs ConfigRestfulResponse
	err = yaml.Unmarshal(body, &configs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body error: %s", err.Error())
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
	exist, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	nodes, _, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes error: %s", err.Error())
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
		return fmt.Errorf("get external controller error: %s", err.Error())
	}

	selector, err := getSelector()
	if err != nil {
		return fmt.Errorf("get selector error: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/proxies/%s", externalController, selector), nil)
	if err != nil {
		return fmt.Errorf("new request error: %s", err.Error())
	}

	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, nodeName)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("change node error: %s", resp.Status)
	}
	return nil
}

func ChangeMode(mode string) error {
	exist, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller error: %s", err.Error())
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("http://%s/configs", externalController), nil)
	if err != nil {
		return fmt.Errorf("new request error: %s", err.Error())
	}

	if mode != "global" && mode != "rule" {
		return fmt.Errorf("unknown mode: %s", mode)
	}
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"mode": "%s"}`, mode)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("change mode error: %s", resp.Status)
	}

	return nil
}

type DelayTestResponse struct {
	Delay   int    `json:"delay"`
	Message string `json:"message"`
}

func BenchmarkNode() error {
	exist, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	externalController, err := getExternalController()
	if err != nil {
		return fmt.Errorf("get external controller error: %s", err.Error())
	}

	nodes, _, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes error: %s", err.Error())
	}

	for _, node := range nodes {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/proxies/%s/delay?timeout=%d&url=%s", externalController, node.Name, 3000, "http://www.gstatic.com/generate_204"), nil)
		if err != nil {
			return fmt.Errorf("new request error: %s", err.Error())
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("do request error: %s", err.Error())
		}
		//goland:noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body error: %s", err.Error())
		}

		var delayTestResponse DelayTestResponse
		err = json.Unmarshal(body, &delayTestResponse)
		if err != nil {
			return fmt.Errorf("unmarshal response body error: %s", err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("*** benchmark node %s, status: %s, message: %s\n", node.Name, resp.Status, delayTestResponse.Message)
			continue
		}

		fmt.Printf("benchmark node %s, status: %s, delay: %d\n", node.Name, resp.Status, delayTestResponse.Delay)
	}

	return nil
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
	exist, _, err := checkUserProcess("cyber-core")
	if err != nil {
		return fmt.Errorf("check user process error: %s", err.Error())
	}

	if !exist {
		return fmt.Errorf("cyber-core is not running")
	}

	nodes, currentNode, err := getNodes()
	if err != nil {
		return fmt.Errorf("get nodes error: %s", err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"当前选择", "节点名称", "节点类型", "延迟", "位置"})

	for _, node := range nodes {
		var data []string
		if node.Name == currentNode {
			data = append(data, "*")
		} else {
			data = append(data, "")
		}
		if node.Delay == 0 {
			data = append(data, node.Name, node.Type, "N/A", getCountryCode(node.Name))
			table.Append(data)
			continue
		}
		data = append(data, node.Name, node.Type, fmt.Sprintf("%dms", node.Delay), getCountryCode(node.Name))
		table.Append(data)
	}

	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Render() // Send output

	return nil
}

func getNodes() ([]Node, string, error) {
	externalController, err := getExternalController()
	if err != nil {
		return nil, "", fmt.Errorf("get external controller error: %s", err.Error())
	}

	selector, err := getSelector()
	if err != nil {
		return nil, "", fmt.Errorf("get selector error: %s", err.Error())
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, selector))
	if err != nil {
		return nil, "", fmt.Errorf("get proxies error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body error: %s", err.Error())
	}

	var proxies ProxySelectorResponse

	err = json.Unmarshal(body, &proxies)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal body error: %s", err.Error())
	}

	var nodes []Node
	for _, nodeName := range proxies.All {
		nodeResp, err := http.Get(fmt.Sprintf("http://%s/proxies/%s", externalController, nodeName))
		if err != nil {
			return nil, "", fmt.Errorf("get proxy error: %s", err.Error())
		}
		//goland:noinspection GoUnhandledErrorResult
		defer nodeResp.Body.Close()

		nodeBody, err := io.ReadAll(nodeResp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("read node body error: %s", err.Error())
		}

		var node ProxyRestfulResponse
		err = json.Unmarshal(nodeBody, &node)
		if err != nil {
			return nil, "", fmt.Errorf("unmarshal node body error: %s", err.Error())
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

func getExternalController() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("get current user error: %s", err.Error())
	}

	configFilePath := filepath.Join(usr.HomeDir, ".cyber", "node", "config.yaml")

	// 打开配置文件，解析yaml
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return "", fmt.Errorf("open config file error: %s", err.Error())
	}
	//goland:noinspection GoUnhandledErrorResult
	defer configFile.Close()

	config := make(map[any]any)
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return "", fmt.Errorf("decode config file error: %s", err.Error())
	}

	externalController, ok := config["external-controller"]
	if !ok {
		return "", fmt.Errorf("external-controller not found")
	}

	return externalController.(string), nil
}

//goland:noinspection GoUnusedParameter
func getCountryCode(nodeName string) string {
	return "N/A"
}
