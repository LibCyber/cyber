/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall cyber core service",
	Long: `Uninstall cyber core service.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.UninstallService()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Uninstall success.")
	},
}

func init() {
	ServiceCmd.AddCommand(uninstallCmd)
}
