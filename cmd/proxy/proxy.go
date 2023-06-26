/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"github.com/LibCyber/cyber/cmd/proxy/unset"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// ProxyCmd represents the proxy command
var ProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Manage proxy",
	Long: `Manage proxy.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	ProxyCmd.AddCommand(unset.UnsetCmd)
}
