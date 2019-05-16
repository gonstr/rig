package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of rig",
	Long:  `All software has versions. This is Rig's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rig v0.3.4")
	},
}
