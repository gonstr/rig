package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rig",
	Short: "Rig is a code-gen tool for k8s manifests",
	Long: `A code-gen tool for k8s manifests.
	
Manage your manifests in versioned templates hosted in any git repository.
Complete documentation is available at https://github.com/gonstr/rig.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
