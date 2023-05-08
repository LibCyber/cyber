/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Setup proxy for shell",
	Long: `Setup proxy for shell.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, socksPort, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Println("Copy and run the following commands:")
		fmt.Printf("export http_proxy=http://127.0.0.1:%d;\nexport https_proxy=http://127.0.0.1:%d;\nexport all_proxy=socks5://127.0.0.1:%d;\n", httpPort, httpPort, socksPort)
	},
}

func init() {
	ProxyCmd.AddCommand(shellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shellCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shellCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
