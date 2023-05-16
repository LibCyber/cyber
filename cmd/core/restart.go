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

// restartCmd represents the start command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart cyber core",
	Long: `Restart cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.Stop()
		if err != nil && !core.IsNotRunning(err) {
			util.PrintlnExit(err)
		}

		pid, err := core.Start()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Printf("Restarted cyber-core with PID %d\n", pid)
	},
}

func init() {
	CoreCmd.AddCommand(restartCmd)
}
