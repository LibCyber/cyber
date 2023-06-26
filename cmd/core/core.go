/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// CoreCmd represents the core command
var CoreCmd = &cobra.Command{
	Use:   "core",
	Short: "Manage core",
	Long: `Manage core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
