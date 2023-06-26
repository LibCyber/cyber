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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install core as a service",
	Long: `Install core as a service so that it can be started automatically.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.InstallService()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Install success, please run `cyber core start` to start cyber core service.")
	},
}

func init() {
	ServiceCmd.AddCommand(installCmd)
}
