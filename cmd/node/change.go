/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change node",
	Long: `Change node.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			util.PrintlnExit("Please specify node")
		}

		node := args[0]
		err := core.ChangeNode(node)
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Node changed to", node)
	},
}

func init() {
	NodeCmd.AddCommand(changeCmd)
}
