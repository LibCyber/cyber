/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package tunnel

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// TunnelCmd represents the tunnel command
var TunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "Start/Stop tunnel (Experimental)",
	Long: `Start/Stop tunnel (Experimental).
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
