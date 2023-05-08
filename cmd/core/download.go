/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download cyber core",
	Long: `Download cyber core.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading cyber core...")
		err := core.Download()
		if err != nil {
			util.PrintlnExit(err.Error())
		}
		fmt.Println("Download cyber core successfully!")
	},
}

func init() {
	CoreCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
