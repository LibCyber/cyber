/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"LibCyber/cyber/internal/core"
	"LibCyber/cyber/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
)

// readmeCmd represents the git command
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Setup proxy for readme",
	Long: `Setup proxy for readme.
`,
	Run: func(cmd *cobra.Command, args []string) {
		httpPort, socksPort, err := core.GetProxyPort()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Currently cyber-core is listening on port", httpPort, "for http proxy, and port", socksPort, "for socks5 proxy.")
		fmt.Println("If you want to use http proxy or socks proxy on telegram, stream, firefox or other applications, follow the steps below to set up proxy manually.")
		fmt.Println()
		fmt.Println("For http proxy:")
		fmt.Println("  1. Open the settings page of the application you want to use proxy on.")
		fmt.Println("  2. Find the proxy settings page.")
		fmt.Println("  3. Select http proxy.")
		fmt.Println("  4. Set the proxy server to 127.0.0.1 and port to", httpPort, ".")
		fmt.Println("  5. Save the settings.")

		fmt.Println()
		fmt.Println("For socks5 proxy:")
		fmt.Println("  1. Open the settings page of the application you want to use proxy on.")
		fmt.Println("  2. Find the proxy settings page.")
		fmt.Println("  3. Select socks5 proxy.")
		fmt.Println("  4. Set the proxy server to 127.0.0.1 and port to", socksPort, ".")
		fmt.Println("  5. Save the settings.")

		fmt.Println()
		fmt.Println("If you want to use proxy on git, npm, yarn, you can run the following commands to set up proxy automatically:")
		fmt.Println("cyber proxy git")
		fmt.Println("cyber proxy npm")
		fmt.Println("cyber proxy yarn")

		fmt.Println()
		fmt.Println("If you want to use proxy on docker, follow the steps below to set up proxy manually.")
		fmt.Println("For Linux:")
		fmt.Println("  1. Find daemon.json file of docker. Usually it is located at /etc/docker/daemon.json.")
		fmt.Println("  2. Add the following lines to daemon.json:")
		fmt.Println("     {")
		fmt.Println("       \"proxies\": {")
		fmt.Println("         \"default\": {")
		fmt.Printf(`           "httpProxy": "http://127.0.0.1:%d",`, httpPort)
		fmt.Println()
		fmt.Printf(`           "httpsProxy": "http://127.0.0.1:%d",`, httpPort)
		fmt.Println()
		fmt.Printf(`           "noProxy": "localhost,127.0.0.0/8"`)
		fmt.Println()
		fmt.Println("         }")
		fmt.Println("       }")
		fmt.Println("     }")
		fmt.Println("  3. Restart docker.")
		fmt.Println()
		fmt.Println("For Windows and macOS:")
		fmt.Println("  1. Open the settings page of docker.")
		fmt.Println("  2. Find the proxy settings page.")
		fmt.Println("  3. Select manual proxy configuration.")
		fmt.Println("  4. Set the proxy server to the following values:")
		fmt.Printf("     http://127.0.0.1:%d for http proxy,\n", httpPort)
		fmt.Printf("     http://127.0.0.1:%d for https proxy,\n", httpPort)
		fmt.Println("     localhost,127.0.0.0/8 for bypass proxy.")
		fmt.Println("  5. Save the settings.")
		fmt.Println("  6. Restart docker.")

	},
}

func init() {
	ProxyCmd.AddCommand(readmeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
