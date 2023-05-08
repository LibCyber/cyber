/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package tunnel

import (
	"fmt"

	"github.com/spf13/cobra"
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Stop tunnel",
	Long: `Stop tunnel.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("off called")
	},
}

func init() {
	TunnelCmd.AddCommand(offCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// offCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// offCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
