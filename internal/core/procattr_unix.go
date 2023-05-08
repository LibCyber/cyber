//go:build !windows

package core

import "syscall"

func procAttrWithNewProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}
