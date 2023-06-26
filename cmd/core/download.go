/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package core

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
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
			util.PrintlnExit(err)
		}
		fmt.Println("Download cyber core successfully!")
	},
}

func init() {
	CoreCmd.AddCommand(downloadCmd)
}
