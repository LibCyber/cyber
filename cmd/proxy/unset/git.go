/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package unset

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Unset proxy for git",
	Long: `Unset proxy for git.
`,
	Run: func(cmd *cobra.Command, args []string) {
		unsetGitProxy("http")
		unsetGitProxy("https")
	},
}

func init() {
	UnsetCmd.AddCommand(gitCmd)
}

func unsetGitProxy(scheme string) {
	cmd := exec.Command("git", "config", "--global", "--unset", fmt.Sprintf("%s.proxy", scheme))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Setting %s proxy: %v\n", scheme, err)
		return
	}
	fmt.Printf("Successfully unset %s proxy\n", scheme)
}
