package api

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"os"
	"testing"
)

func TestGetProductInfo(t *testing.T) {
	platform, _, platformVersion, _ := host.PlatformInformation()

	fmt.Println("OS Info: ", platform, platformVersion)

	hostname, _ := os.Hostname()
	fmt.Println("Model Name: ", hostname)

	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) == 0 {
		fmt.Println("Getting CPU info")
		return
	}
	fmt.Println("Model: ", cpuInfo[0].ModelName)
}
