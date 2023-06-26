/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cyber core status",
	Long: `Show cyber core status.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := core.Status()
		if err != nil {
			util.PrintlnExit(err)
		}
		if pid > 0 {
			fmt.Println("cyber core is running, pid:", pid)
			return
		}

		fmt.Println("cyber core is not running")
	},
}

func init() {
	CoreCmd.AddCommand(statusCmd)
}
