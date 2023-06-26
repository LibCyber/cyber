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

// pnpmCmd represents the git command
var pnpmCmd = &cobra.Command{
	Use:   "pnpm",
	Short: "Unset proxy for pnpm",
	Long: `Unset proxy for pnpm.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := unsetPnpmProxy("http")
		if err != nil {
			util.PrintlnExit(err)
		}
		err = unsetPnpmProxy("https")
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Unset pnpm proxy done.")
	},
}

func init() {
	UnsetCmd.AddCommand(pnpmCmd)
}

func unsetPnpmProxy(scheme string) error {
	var cmd *exec.Cmd
	if scheme == "http" {
		cmd = exec.Command("pnpm", "config", "delete", "proxy")
	} else if scheme == "https" {
		cmd = exec.Command("pnpm", "config", "delete", "https-proxy")
	} else {
		return fmt.Errorf("scheme must be http or https")
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("setting %s proxy: %v", scheme, err)
	}

	return nil
}
