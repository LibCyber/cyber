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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
