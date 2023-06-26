/*
Copyright © 2023 LibCyber Team libcyberstudio@gmail.com
*/
package node

import (
	"fmt"
	"github.com/LibCyber/cyber/internal/app"
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
		if app.Language() == "zh" {
			fmt.Println("开始测试延迟...")
		} else {
			fmt.Println("Start benchmarking...")
		}
		err := core.BenchmarkNode()
		if err != nil {
			util.PrintlnExit(err)
		}

		if app.Language() == "zh" {
			fmt.Println("延迟测试完成。")
		} else {
			fmt.Println("Benchmark finished.")
		}
	},
}

func init() {
	NodeCmd.AddCommand(benchmarkCmd)
}
