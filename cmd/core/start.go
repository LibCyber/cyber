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
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cyber core",
	Long: `Start cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := core.Start()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Printf("Started cyber core with PID %d\n", pid)
	},
}

func init() {
	CoreCmd.AddCommand(startCmd)
}
