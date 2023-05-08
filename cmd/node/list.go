/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
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
			util.PrintlnExit(err.Error())
		}
	},
}

func init() {
	NodeCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
