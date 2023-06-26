/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package secret

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off API endpoint secret",
	Long: `Turn off API endpoint secret.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.ModifyAPISecret(false, "")
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	SecretCmd.AddCommand(offCmd)
}
