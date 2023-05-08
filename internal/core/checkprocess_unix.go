//go:build !windows
// +build !windows

package core

import (
	"github.com/shirou/gopsutil/v3/process"
	"os/user"
	"strconv"
	"strings"
)

// check if the process is running, if so, return true, pid, name, nil
func checkUserProcess(processName string) (bool, int, string, error) {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		return false, 0, "", err
	}

	// 将当前用户的UID转换为整数
	uid, err := strconv.Atoi(currentUser.Uid)
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
			uids, err := p.Uids()
			if err != nil || len(uids) == 0 {
				continue
			}

			if uids[0] == int32(uid) {
				return true, int(p.Pid), name, nil
			}
		}
	}

	return false, 0, "", nil
}
