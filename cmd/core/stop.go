/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop cyber core",
	Long: `Stop cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.Stop()
		if err != nil {
			util.PrintlnExit(err.Error())
		}
	},
}

func init() {
	CoreCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
