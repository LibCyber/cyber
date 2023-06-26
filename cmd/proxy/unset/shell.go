/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package unset

import (
	"fmt"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Unset proxy for shell",
	Long: `Unset proxy for shell.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Copy and run the following commands:")
		fmt.Println("unset http_proxy;\nunset https_proxy;\nunset all_proxy;")
	},
}

func init() {
	UnsetCmd.AddCommand(shellCmd)
}
