/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/core"
	"github.com/LibCyber/cyber/pkg/util"
	"github.com/spf13/cobra"
)

// benchmarkCmd represents the benchmark command
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run batch latency testing on all nodes",
	Long: `Run batch latency testing on all nodes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start benchmarking...")
		err := core.BenchmarkNode()
		if err != nil {
			util.PrintlnExit(err)
		}

		fmt.Println("Benchmark finished.")
	},
}

func init() {
	NodeCmd.AddCommand(benchmarkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benchmarkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchmarkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
