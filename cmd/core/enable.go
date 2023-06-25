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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
