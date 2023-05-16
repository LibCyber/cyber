/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show user information",
	Long: `Show user information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var c = api.NewClient()
		if app.Language() == "zh" {
			fmt.Println("正在获取用户信息...")
		} else {
			fmt.Println("Getting user information...")
		}
		info, err := c.GetUserInfo()
		if err != nil {
			util.PrintlnExit(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		if app.Language() == "zh" {
			table.Append([]string{"用户名", info.Account})
			table.Append([]string{"已用流量", fmt.Sprintf("%.2fGB", float64(info.U+info.D)/(1024*1024*1024))})
			table.Append([]string{"总流量", fmt.Sprintf("%.2fGB", float64(info.TransferEnable)/(1024*1024*1024))})
			table.Append([]string{"余额", fmt.Sprintf("￥%.2f", float64(info.Balance)/100)})
			table.Append([]string{"过期时间", info.ExpiredAt})
			table.Append([]string{"最后登录时间", time.Unix(int64(info.LastLogin), 0).Format("2006-01-02 15:04:05")})
			table.Append([]string{"流量重置日", fmt.Sprintf("每月%d号", info.ResetTime)})
			table.Append([]string{"账户状态", formatStatus(info.Status)})
			table.Append([]string{"节点状态", formatStatus(info.Enable)})
		} else {
			table.Append([]string{"Username", info.Account})
			table.Append([]string{"Used Traffic", fmt.Sprintf("%.2fGB", float64(info.U+info.D)/(1024*1024*1024))})
			table.Append([]string{"Total Traffic", fmt.Sprintf("%.2fGB", float64(info.TransferEnable)/(1024*1024*1024))})
			table.Append([]string{"Balance", fmt.Sprintf("￥%.2f", float64(info.Balance)/100)})
			table.Append([]string{"Expire Time", info.ExpiredAt})
			table.Append([]string{"Last Login Time", time.Unix(int64(info.LastLogin), 0).Format("2006-01-02 15:04:05")})
			table.Append([]string{"Traffic Reset Day", fmt.Sprintf("the %d-th day of every month", info.ResetTime)})
			table.Append([]string{"Account Status", formatStatus(info.Status)})
			table.Append([]string{"Node Enable", formatStatus(info.Enable)})
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func formatStatus(status int) string {
	if status == 0 {
		if app.Language() == "zh" {
			return "禁用"
		}
		return "Disabled"
	} else if status == 1 {
		if app.Language() == "zh" {
			return "启用"
		}
		return "Enabled"
	}

	return "Unknown"
}
