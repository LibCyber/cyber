/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Setup proxy for git",
	Long: `Setup proxy for git.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, _, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}

		setupGitProxy("http", fmt.Sprintf("http://127.0.0.1:%d", httpPort))
		setupGitProxy("https", fmt.Sprintf("http://127.0.0.1:%d", httpPort))

		fmt.Println("git proxy setup")
	},
}

func init() {
	ProxyCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupGitProxy(scheme, proxyUrl string) {
	cmd := exec.Command("git", "config", "--global", fmt.Sprintf("%s.proxy", scheme), proxyUrl)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error setting %s proxy: %v\n", scheme, err)
		return
	}
	fmt.Printf("Successfully set %s proxy to %s\n", scheme, proxyUrl)
}
