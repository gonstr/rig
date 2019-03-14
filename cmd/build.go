package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a rig.yaml template to stdout",
	Long: `Builds a rig.yaml template to stdout. Values are read from
rig.yaml and from command line args.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		check(err)

		filePath := path.Join(wd, "rig.yaml")

		tmpl, err := template.NewFromFile(filePath)
		check(err)

		err = tmpl.Sync()
		check(err)

		yaml, err := tmpl.Build(filePath)
		check(err)

		fmt.Print(yaml)
	},
}
