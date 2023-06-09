/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package core

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable cyber core auto-start on boot",
	Long: `Enable cyber core auto-start on boot.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.EnableService()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	CoreCmd.AddCommand(enableCmd)
}
