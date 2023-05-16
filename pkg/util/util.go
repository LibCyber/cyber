package util

import (
	"fmt"
	"os"
)

func PrintlnExit(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(2)
}
