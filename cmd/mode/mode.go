/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package mode

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// ModeCmd represents the mode command
var ModeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Manage cyber routing mode",
	Long: `Manage cyber routing mode.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
