/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all data",
	Long: `Remove all data of cyber.
`,
	Run: func(cmd *cobra.Command, args []string) {
		isForce, err := cmd.Flags().GetBool("force")
		if err != nil {
			util.PrintlnExit(err)
		}

		if !isForce {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Are you sure to purge all data? (y/N): ")
			text, err := reader.ReadString('\n')
			if err != nil {
				util.PrintlnExit(err)
			}
			text = strings.TrimSpace(text)
			if text != "y" && text != "Y" {
				return
			}
		}

		// stop core
		fmt.Println("Stopping cyber core...")
		err = core.Stop()
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

	purgeCmd.Flags().BoolP("force", "f", false, "Force purge without prompt")
}
