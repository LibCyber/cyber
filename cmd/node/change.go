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
		if len(args) != 1 {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
