package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func PrintlnExit(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(2)
}

func IsAdmin() (bool, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("NET", "SESSION")
	case "linux", "darwin":
		cmd = exec.Command("id", "-u")
	default:
		return false, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, err
	}

	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		return true, nil
	} else {
		// For Unix-like systems, the root user has an ID of 0.
		return strings.TrimSpace(string(output)) == "0", nil
	}
}

func CopyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath) //打开源文件
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer srcFile.Close()

	destFile, err := os.Create(destPath) //创建目标文件
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) //复制文件
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

//func Pointer2Bool(p bool) *bool {
//	return &p
//}
