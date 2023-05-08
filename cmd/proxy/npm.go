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

// npmCmd represents the git command
var npmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Setup proxy for npm",
	Long: `Setup proxy for npm.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, _, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}

		err = setupNpmProxy("http", fmt.Sprintf("http://127.0.0.1:%d", httpPort))
		if err != nil {
			util.PrintlnExit(err)
		}
		err = setupNpmProxy("https", fmt.Sprintf("http://127.0.0.1:%d", httpPort))
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Setup npm proxy done.")
	},
}

func init() {
	ProxyCmd.AddCommand(npmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupNpmProxy(scheme, proxyUrl string) error {
	var proxyText string
	if scheme == "http" {
		proxyText = "proxy"
	} else if scheme == "https" {
		proxyText = "https-proxy"
	} else {
		return fmt.Errorf("scheme must be http or https")
	}
	cmd := exec.Command("npm", "config", "set", proxyText, proxyUrl)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting %s proxy: %v", scheme, err)
	}

	return nil
}
