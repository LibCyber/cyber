/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all data",
	Long: `Remove all data of cyber.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// stop core
		fmt.Println("Stopping cyber core...")
		err := core.Stop()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Logging out...")
		var c = api.NewClient()
		err = c.Logout()
		if err != nil {
			println("Logout failed:", err.Error())
		} else {
			fmt.Println("Logout successfully")
		}

		// remove binaries and data
		fmt.Println("Removing binaries and data...")
		err = core.Purge()
		if err != nil {
			util.PrintlnExit(err)
		}

		// remove service
		if isServiceMode, err := core.IsServiceInstalled(); err != nil && isServiceMode {
			err = core.UninstallService()
			if err != nil {
				util.PrintlnExit(err)
			}
		}

		fmt.Println("Purge done")
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// purgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
