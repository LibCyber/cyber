/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop cyber core",
	Long: `Stop cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.Stop()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	CoreCmd.AddCommand(stopCmd)
}
