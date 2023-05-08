/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package unset

import (
	"fmt"
	"github.com/LibCyber/cyber/pkg/util"
	"os/exec"

	"github.com/spf13/cobra"
)

// npmCmd represents the git command
var npmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Unset proxy for npm",
	Long: `Unset proxy for npm.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := unsetNpmProxy("http")
		if err != nil {
			util.PrintlnExit(err)
		}
		err = unsetNpmProxy("https")
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Unset npm proxy done.")
	},
}

func init() {
	UnsetCmd.AddCommand(npmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func unsetNpmProxy(scheme string) error {
	var cmd *exec.Cmd
	if scheme == "http" {
		cmd = exec.Command("npm", "config", "delete", "proxy")
	} else if scheme == "https" {
		cmd = exec.Command("npm", "config", "delete", "https-proxy")
	} else {
		return fmt.Errorf("scheme must be http or https")
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting %s proxy: %v", scheme, err)
	}

	return nil
}
