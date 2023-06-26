/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// NodeCmd represents the node command
var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage nodes",
	Long: `Manage nodes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
