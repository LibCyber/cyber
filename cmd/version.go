/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/constant"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long: `Print the version of the application.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cyber", "version", constant.APP_VERSION, fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		// TODO: check update

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
