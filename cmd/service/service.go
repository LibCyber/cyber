/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// ServiceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage service",
	Long: `Manage service.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			util.PrintlnExit(err)
		}
	},
}

func init() {}
