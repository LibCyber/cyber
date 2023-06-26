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
		exists, err := c.IsNodeExists()
		if err != nil {
			util.PrintlnExit(err)
		}
		if exists {
			util.PrintlnExit("Node already exists, please use `cyber node update` to update nodes.")
		}

		fmt.Println("Downloading nodes ...")
		err = c.DownloadNodes()
		if err != nil {
			util.PrintlnExit(err)
		}
		fmt.Println("Download nodes successfully.")
	},
}

func init() {
	NodeCmd.AddCommand(downloadCmd)
}
