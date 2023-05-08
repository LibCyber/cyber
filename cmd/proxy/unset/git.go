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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func unsetGitProxy(scheme string) {
	cmd := exec.Command("git", "config", "--global", "--unset", fmt.Sprintf("%s.proxy", scheme))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error setting %s proxy: %v\n", scheme, err)
		return
	}
	fmt.Printf("Successfully unset %s proxy\n", scheme)
}
