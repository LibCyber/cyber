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
		fmt.Printf("http-port: %d\nsocks-port: %d\n", httpPort, socksPort)
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
