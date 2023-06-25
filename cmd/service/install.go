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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
