/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cyber core status",
	Long: `Show cyber core status.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.Status()
		if err != nil {
			util.PrintlnExit(err.Error())
		}
	},
}

func init() {
	CoreCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
