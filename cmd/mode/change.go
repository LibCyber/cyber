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

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change cyber routing mode",
	Long: `Change cyber routing mode. For example:

cyber mode change <mode>

where <mode> is one of the following:

- rule
- global


`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			util.PrintlnExit("Please specify mode, either 'rule' or 'global'")
		}

		mode := args[0]
		if mode != "rule" && mode != "global" {
			util.PrintlnExit("Invalid mode, please specify either 'rule' or 'global'")
		}

		err := core.ChangeMode(mode)
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Mode changed to", mode)
	},
}

func init() {
	ModeCmd.AddCommand(changeCmd)
}
