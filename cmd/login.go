/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/api"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to LibCyber",
	Long: `Login to LibCyber.
`,
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			util.PrintlnExit(err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			util.PrintlnExit(err)
		}

		if username != "" && password == "" {
			util.PrintlnExit("password is required")
		}

		if username == "" && password == "" {
			// 请求输入用户名和密码
			fmt.Print("Username: ")
			scanln, err := fmt.Scanln(&username)
			if err != nil {
				util.PrintlnExit(err)
			}
			if scanln == 0 {
				util.PrintlnExit("username is required")
			}

			fmt.Print("Password: ")
			scanln, err = fmt.Scanln(&password)
			if err != nil {
				util.PrintlnExit(err)
			}
			if scanln == 0 {
				util.PrintlnExit("password is required")
			}
		}

		fmt.Println()
		fmt.Printf("Logging in as: %s\npassword: %s\n", username, formatStarWords(password))
		var c = api.NewClient()
		_, err = c.Login(username, password)
		if err != nil {
			util.PrintlnExit(err)
		}

		// 成功登录
		fmt.Println()
		fmt.Printf("Successfully logged in as: %s\n", username)
	},
}

func formatStarWords(s string) string {
	var starWords = ""
	for i := 0; i < len(s); i++ {
		starWords += "*"
	}
	return starWords
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username or email used to sign up.")
	loginCmd.Flags().StringP("password", "p", "", "Password.")
}
