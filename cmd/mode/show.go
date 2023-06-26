/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package mode

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cyber routing mode",
	Long: `Show cyber routing mode.
`,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := core.GetConfigs()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Current mode:", configs.Mode)
	},
}

func init() {
	ModeCmd.AddCommand(showCmd)
}
