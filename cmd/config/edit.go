/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit cyber core and nodes config",
	Long: `Edit cyber core and nodes config.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core.EditNodeConfigWithVi()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(editCmd)
}
