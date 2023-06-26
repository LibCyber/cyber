/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List nodes",
	Long: `List nodes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.ListNodes()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	NodeCmd.AddCommand(listCmd)
}
