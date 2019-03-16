package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rig",
	Short: "Rig is a Kubernetes manifest preprocessor and templating tool",
	Long: `Rig is a Kubernetes manifest preprocessor and templating tool.

Complete documentation is available at https://github.com/gonstr/rig.
`,
}

func Execute() {
	err := rootCmd.Execute()
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
