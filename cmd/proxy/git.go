/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
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

		fmt.Println("Git proxy setup successfully")
	},
}

func init() {
	ProxyCmd.AddCommand(gitCmd)
}

func setupGitProxy(scheme, proxyUrl string) {
	cmd := exec.Command("git", "config", "--global", fmt.Sprintf("%s.proxy", scheme), proxyUrl)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Setting %s proxy: %v\n", scheme, err)
		return
	}
	fmt.Printf("Successfully set %s proxy to %s\n", scheme, proxyUrl)
}
