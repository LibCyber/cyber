/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package proxy

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
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

		if app.Language() == "zh" {
			fmt.Println("当前 cyber内核 正在端口", httpPort, "监听 http 代理，以及在端口", socksPort, "监听 socks5 代理。")
			fmt.Println("如果您希望在 python、telegram、steam、firefox 或其他应用程序中使用 http 代理或 socks 代理，请按照以下步骤手动设置代理。")
			fmt.Println()

			fmt.Println("对于 http 代理：")
			fmt.Println("  1. 打开您想要使用代理的应用程序的设置页面。")
			fmt.Println("  2. 找到代理设置页面。")
			fmt.Println("  3. 选择 http 代理。")
			fmt.Println("  4. 将代理服务器设置为 127.0.0.1，端口设置为", httpPort, "。")
			fmt.Println("  5. 保存设置。")

			fmt.Println()
			fmt.Println("对于 socks5 代理：")
			fmt.Println("  1. 打开您想要使用代理的应用程序的设置页面。")
			fmt.Println("  2. 找到代理设置页面。")
			fmt.Println("  3. 选择 socks5 代理。")
			fmt.Println("  4. 将代理服务器设置为 127.0.0.1，端口设置为", socksPort, "。")
			fmt.Println("  5. 保存设置。")

			fmt.Println()
			fmt.Println("如果您想在 git、npm、yarn 上使用代理，您可以运行以下命令自动设置代理：")
			fmt.Println("cyber proxy git")
			fmt.Println("cyber proxy npm")
			fmt.Println("cyber proxy yarn")

			fmt.Println()
			fmt.Println("如果您想在 docker 上使用代理，请按照以下步骤手动设置代理。")
			fmt.Println("对于 Linux：")
			fmt.Println("  1. 找到 docker 的 daemon.json 文件。通常位于 /etc/docker/daemon.json。")
			fmt.Println("  2. 向 daemon.json 添加以下行：")
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
			fmt.Println("  3. 重启 docker。")
			fmt.Println()

			fmt.Println("对于 Windows 和 macOS：")
			fmt.Println("  1. 打开 Docker 的设置页面。")
			fmt.Println("  2. 找到代理设置页面。")
			fmt.Println("  3. 选择手动代理配置。")
			fmt.Println("  4. 将代理服务器设置为以下值：")
			fmt.Printf("     http://127.0.0.1:%d 作为 http 代理,\n", httpPort)
			fmt.Printf("     http://127.0.0.1:%d 作为 https 代理,\n", httpPort)
			fmt.Println("     localhost,127.0.0.0/8 用于绕过代理。")
			fmt.Println("  5. 保存设置。")
			fmt.Println("  6. 重启 Docker。")
			fmt.Println()
			fmt.Println()
			fmt.Println("如果你不想这么麻烦，你可以运行以下命令自动设置全局TUN代理（实验性）：")
			fmt.Println("cyber tunnel on")
			fmt.Println("记住，如果出现网络问题，运行以下命令关闭TUN代理：")
			fmt.Println("cyber tunnel off")

		} else {
			fmt.Println("Currently cyber core is listening on port", httpPort, "for http proxy, and port", socksPort, "for socks5 proxy.")
			fmt.Println("If you want to use http proxy or socks proxy on python, telegram, steam, firefox or other applications, follow the steps below to set up proxy manually.")
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
			fmt.Println()
			fmt.Println()
			fmt.Println("If you just want to keep it simple, you can run the following command to automatically set the global TUN proxy (experimental):")
			fmt.Println("cyber tunnel on")
			fmt.Println("Remember, in case of network problems, run the following command to turn off the TUN proxy:")
			fmt.Println("cyber tunnel off")
		}

	},
}

func init() {
	ProxyCmd.AddCommand(readmeCmd)
}
