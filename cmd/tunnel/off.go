/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package tunnel

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Stop tunnel",
	Long: `Stop tunnel.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.DisableTun()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Tunnel disabled")
	},
}

func init() {
	TunnelCmd.AddCommand(offCmd)
}
