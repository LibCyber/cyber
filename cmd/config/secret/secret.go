/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package secret

import (
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// SecretCmd represents the secret command
var SecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage API endpoint secret",
	Long: `Manage API endpoint secret.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
