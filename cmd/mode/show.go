/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package mode

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cyber routing mode",
	Long: `Show cyber routing mode.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("show called")
		configs, err := core.GetConfigs()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Current mode:", configs.Mode)
	},
}

func init() {
	ModeCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
