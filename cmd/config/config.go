/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"github.com/LibCyber/cyber/cmd/config/secret"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cyber core and node config",
	Long: `Manage cyber core and node config.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(secret.SecretCmd)
}
