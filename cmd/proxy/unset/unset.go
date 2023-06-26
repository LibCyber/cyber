/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package unset

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// UnsetCmd represents the unset command
var UnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset proxy",
	Long: `Unset proxy.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
