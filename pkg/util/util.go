package util

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

func MakeRandStr(length int, isNumber bool) string {
	var chars string
	// 密码字符集，可任意添加你需要的字符
	if isNumber {
		chars = "0123456789"
	} else {
		//chars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	char := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		char += string(chars[rand.Intn(len(chars))])
	}

	return char
}
