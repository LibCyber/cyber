/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
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
		if app.Language() == "zh" {
			fmt.Println("复制并运行以下命令：")
		} else {
			fmt.Println("Copy and run the following commands:")
		}
		fmt.Printf("export http_proxy=http://127.0.0.1:%d;\nexport https_proxy=http://127.0.0.1:%d;\nexport all_proxy=socks5://127.0.0.1:%d;\n", httpPort, httpPort, socksPort)
	},
}

func init() {
	ProxyCmd.AddCommand(shellCmd)
}
