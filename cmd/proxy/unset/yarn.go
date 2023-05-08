/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package unset

import (
	"LibCyber/cyber/pkg/util"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// yarnCmd represents the git command
var yarnCmd = &cobra.Command{
	Use:   "yarn",
	Short: "Unset proxy for yarn",
	Long: `Unset proxy for yarn.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := unsetYarnProxy("http")
		if err != nil {
			util.PrintlnExit(err)
		}
		err = unsetYarnProxy("https")
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Unset yarn proxy done.")
	},
}

func init() {
	UnsetCmd.AddCommand(yarnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func unsetYarnProxy(scheme string) error {
	var cmd *exec.Cmd
	if scheme == "http" {
		cmd = exec.Command("yarn", "config", "delete", "proxy")
	} else if scheme == "https" {
		cmd = exec.Command("yarn", "config", "delete", "https-proxy")
	} else {
		return fmt.Errorf("scheme must be http or https")
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting %s proxy: %v", scheme, err)
	}

	return nil
}
