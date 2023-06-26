/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// languageCmd represents the language command
var languageCmd = &cobra.Command{
	Use:   "language",
	Short: "Mange language",
	Long: `Mange language, currently support:
	- en
	- zh (partially translated)
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if app.Language() == "zh" {
				fmt.Println("请指定一种语言，目前支持：en, zh")
				fmt.Println("示例：cyber language en")
				return
			} else {
				fmt.Println("Please specify a language, currently support: en, zh")
				fmt.Println("Example: cyber language en")
				return
			}
		}

		switch args[0] {
		case "en":
			err := app.SetLanguage("en")
			if err != nil {
				util.PrintlnExit(err)
			}
		case "zh":
			err := app.SetLanguage("zh")
			if err != nil {
				util.PrintlnExit(err)
			}
		default:
			err := cmd.Help()
			if err != nil {
				util.PrintlnExit(err)
			}
		}

		if app.Language() == "zh" {
			fmt.Println("语言设置成功")
		} else {
			fmt.Println("Language set successfully")
		}

	},
}

func init() {
	rootCmd.AddCommand(languageCmd)
}
