package sysinfo

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"os"
)

func GetOSInfo() string {
	platform, _, platformVersion, err := host.PlatformInformation()
	if err != nil {
		return "Unknown"
	}
	return platform + " " + platformVersion
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Unknown"
	}
	return hostname
}

func GetModel() string {
	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) == 0 {
		return "Unknown"
	}
	return cpuInfo[0].ModelName
}
