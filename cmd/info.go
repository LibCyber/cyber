/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
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
		fmt.Println("Getting user information...")
		info, err := c.GetUserInfo()
		if err != nil {
			util.PrintlnExit(err.Error())
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{"Username", info.Account})
		table.Append([]string{"Used Traffic", fmt.Sprintf("%.2fGB", float64(info.U+info.D)/(1024*1024*1024))})
		table.Append([]string{"Total Traffic", fmt.Sprintf("%.2fGB", float64(info.TransferEnable)/(1024*1024*1024))})
		table.Append([]string{"Balance", fmt.Sprintf("￥%.2f", float64(info.Balance)/100)})
		table.Append([]string{"Expire Time", info.ExpiredAt})
		table.Append([]string{"Last Login Time", time.Unix(int64(info.LastLogin), 0).Format("2006-01-02 15:04:05")})
		table.Append([]string{"Traffic Reset Day", fmt.Sprintf("the %d-th day of every month", info.ResetTime)})
		table.Append([]string{"Account Status", formatStatus(info.Status)})
		table.Append([]string{"Node Enable", formatStatus(info.Enable)})
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
		return "Disabled"
	} else if status == 1 {
		return "Enabled"
	}

	return "Unknown"
}
