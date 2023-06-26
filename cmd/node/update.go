/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package node

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update nodes",
	Long: `Update nodes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var c = api.NewClient()
		exists, err := c.IsNodeExists()
		if err != nil {
			util.PrintlnExit(err)
		}
		if !exists {
			util.PrintlnExit("Node does not exist, please use `cyber node download` to download nodes.")
		}

		fmt.Println("Updating nodes.")
		err = c.UpdateNodes()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Update nodes successfully.")
	},
}

func init() {
	NodeCmd.AddCommand(updateCmd)
}
