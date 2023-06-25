/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package tunnel

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Start tunnel",
	Long: `Start tunnel.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if app.Language() == "zh" {
			fmt.Println("警告：Tun功能仍处于实验阶段。如果开启Tun后出现网络问题，您可以使用 `cyber tunnel off` 关闭此功能。")
		} else {
			fmt.Println("Warning: The tunnel function is still in the experimental stage. If it causes network issues, you can turn off this function using `cyber tunnel off`.")
		}
		err := core.EnableTun()
		if err != nil {
			util.PrintlnExit(err)
		}

		if app.Language() == "zh" {
			fmt.Println("Tun 启用成功")
		} else {
			fmt.Println("Tunnel enabled")
		}
	},
}

func init() {
	TunnelCmd.AddCommand(onCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
