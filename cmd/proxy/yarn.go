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

// yarnCmd represents the git command
var yarnCmd = &cobra.Command{
	Use:   "yarn",
	Short: "Setup proxy for yarn",
	Long: `Setup proxy for yarn.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, _, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}

		err = setupYarnProxy("http", fmt.Sprintf("http://127.0.0.1:%d", httpPort))
		if err != nil {
			util.PrintlnExit(err)
		}
		err = setupYarnProxy("https", fmt.Sprintf("http://127.0.0.1:%d", httpPort))
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Setup yarn proxy done.")
	},
}

func init() {
	ProxyCmd.AddCommand(yarnCmd)
}

func setupYarnProxy(scheme, proxyUrl string) error {
	var proxyText string
	if scheme == "http" {
		proxyText = "proxy"
	} else if scheme == "https" {
		proxyText = "https-proxy"
	} else {
		return fmt.Errorf("scheme must be http or https")
	}
	cmd := exec.Command("yarn", "config", "set", proxyText, proxyUrl)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("setting %s proxy: %s", scheme, err.Error())
	}

	return nil
}
