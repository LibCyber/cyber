/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from LibCyber",
	Long: `Logout from LibCyber.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping cyber-core...")
		err := core.Stop()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Logging out...")
		var c = api.NewClient()
		err = c.Logout()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Logout successfully")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
