/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if cyber is healthy",
	Long: `Check if cyber and cyber core are installed correctly, if not, try to fix it.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented yet")

		// check if logged in (if access token exists and active)
		c := api.NewClient()
		if c.IsTokenValid() {
			fmt.Println("token is valid")
		} else {
			fmt.Println("token is invalid")
		}
		// check if cyber core is installed
		installed, err := core.IsCoreInstalled()
		if err != nil {
			fmt.Println(err)
		} else {
			if installed {
				fmt.Println("cyber core is installed")
			} else {
				fmt.Println("cyber core is not installed")
			}
		}

		// check if node config exists
		exists, err := c.IsNodeExists()
		if err != nil {
			fmt.Println(err)
		} else {
			if exists {
				fmt.Println("node config exists")
			} else {
				fmt.Println("node config does not exist")
			}
		}

		// check if cyber core is installed as a service
		serviceInstalled, err := core.IsServiceInstalled()
		if err != nil {
			fmt.Println(err)
		} else {
			if serviceInstalled {
				fmt.Println("cyber core is installed as a service")
			} else {
				fmt.Println("cyber core is not installed as a service")
			}
		}

		// check if tunnel on
		enabled, err := core.IsTunEnabled()
		if err != nil {
			fmt.Println(err)
		} else {
			if enabled {
				fmt.Println("tunnel is enabled")
			} else {
				fmt.Println("tunnel is not enabled")
			}
		}

		// check if network under cyber core proxy is available
		// send request to ipforge.libcyber.xyz under proxy and check if response is 200
		proxyPort, _, err := core.GetProxyPort()
		if err != nil {
			fmt.Println(err)
		} else {
			//fmt.Println("proxy port: ", proxyPort)
			ip, err := getIP(proxyPort)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("network under cyber core proxy is available, ip: ", ip)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doctorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doctorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getIP(proxyPort int) (string, error) {
	req, err := http.NewRequest("GET", "http://ipforge.libcyber.xyz", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("127.0.0.1:%d", proxyPort),
			}),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return string(body), nil
	}

	return "", fmt.Errorf("abnormal status code: %d", resp.StatusCode)
}
