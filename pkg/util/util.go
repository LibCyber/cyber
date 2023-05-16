package util

import (
	"fmt"
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
