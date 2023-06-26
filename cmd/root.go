/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"github.com/LibCyber/cyber/cmd/config"
	"github.com/LibCyber/cyber/cmd/core"
	"github.com/LibCyber/cyber/cmd/mode"
	"github.com/LibCyber/cyber/cmd/node"
	"github.com/LibCyber/cyber/cmd/proxy"
	"github.com/LibCyber/cyber/cmd/service"
	"github.com/LibCyber/cyber/cmd/tunnel"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cyber",
	Short: "LibCyber CLI",
	Long: `cyber is a command line client for LibCyber.

TLDR, only 5 steps to liberate your cyber:
  cyber login                1. Login to LibCyber.
  cyber init                 2. Init cyber.
  cyber proxy                3. Setup proxy for git, npm, etc. 

Find more information at: https://docs.libcyber.org/docs/quan-xin-libcyber-ke-hu-duan/cli
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(node.NodeCmd)
	rootCmd.AddCommand(proxy.ProxyCmd)
	rootCmd.AddCommand(tunnel.TunnelCmd)
	rootCmd.AddCommand(core.CoreCmd)
	rootCmd.AddCommand(mode.ModeCmd)
	rootCmd.AddCommand(service.ServiceCmd)
	rootCmd.AddCommand(config.ConfigCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cyber.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
