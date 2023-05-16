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

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show proxy information",
	Long: `Show proxy information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, socksPort, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}
		if app.Language() == "zh" {
			fmt.Printf("当前 cyber-core 正在本机127.0.0.1，端口 %d 监听 http 代理，以及在端口 %d 监听 socks5 代理。\n", httpPort, socksPort)
		} else {
			fmt.Printf("The cyber-core is listening http proxy on port %d and socks5 proxy on port %d.\n", httpPort, socksPort)
		}
	},
}

func init() {
	ProxyCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
