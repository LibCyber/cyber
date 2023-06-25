/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package core

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable cyber core auto-start on boot",
	Long: `Disable cyber core auto-start on boot.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.DisableService()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	CoreCmd.AddCommand(disableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
