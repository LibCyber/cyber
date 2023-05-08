/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download nodes",
	Long: `Download nodes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var c = api.NewClient()
		fmt.Println("Downloading nodes ...")
		err := c.DownloadNodes()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Println("Download nodes successfully.")
	},
}

func init() {
	NodeCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
