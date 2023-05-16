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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ProxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ProxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
