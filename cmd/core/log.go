/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package core

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show log of cyber core",
	Long: `Show log of cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		isFollow, err := cmd.Flags().GetBool("follow")
		if err != nil {
			util.PrintlnExit(err)
		}

		if isFollow {
			err := core.ShowLogFollow()
			if err != nil {
				util.PrintlnExit(err)
			}
		} else {
			err := core.ShowLog()
			if err != nil {
				util.PrintlnExit(err)
			}
		}
	},
}

func init() {
	CoreCmd.AddCommand(logCmd)

	// -f 添加此标志，将不会退出，而是一直打印日志
	logCmd.Flags().BoolP("follow", "f", false, "Follow log")
}
