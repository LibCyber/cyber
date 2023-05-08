//go:build windows
// +build windows

package core

import "syscall"

const CREATE_NO_WINDOW = 0x08000000

func procAttrWithNewProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW,
	}
}
