//go:build windows
// +build windows

package core

import (
	"github.com/shirou/gopsutil/v3/process"

	"os/user"
	"strings"
)

// check if the process is running, if so, return true, pid, name, nil
func checkUserProcess(processName string) (bool, int, string, error) {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		return false, 0, "", err
	}

	// 获取系统中所有运行中的进程
	processes, err := process.Processes()
	if err != nil {
		return false, 0, "", err
	}

	// 遍历所有进程，检查进程名和用户ID
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		if strings.Contains(name, processName) {
			username, err := p.Username()
			if err != nil {
				continue
			}

			if username == currentUser.Username {
				return true, int(p.Pid), name, nil
			}
		}
	}

	return false, 0, "", nil
}
