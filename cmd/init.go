/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cyber",
	Long: `Initialize cyber, including:
- download cyber core
- download nodes
- start cyber core
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading cyber core...")
		err := core.Download()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Println("Download cyber core successfully!")

		var c = api.NewClient()
		fmt.Println("Downloading nodes ...")
		err = c.DownloadNodes()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Println("Download nodes successfully.")

		pid, err := core.Start()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Printf("Started cyber-core with PID %d\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
