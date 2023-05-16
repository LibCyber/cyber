/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package mode

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// ModeCmd represents the mode command
var ModeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Manage cyber routing mode",
	Long: `Manage cyber routing mode.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
