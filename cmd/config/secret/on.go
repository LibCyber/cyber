/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package secret

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"

	"github.com/spf13/cobra"
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on API endpoint secret",
	Long: `Turn on API endpoint secret.
`,
	Run: func(cmd *cobra.Command, args []string) {
		secret := util.MakeRandStr(16, false)
		err := core.ModifyAPISecret(true, secret)
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("API endpoint secret is:", secret)
	},
}

func init() {
	SecretCmd.AddCommand(onCmd)
}
