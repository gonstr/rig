package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

var values []string
var stringValues []string

func init() {
	buildCmd.Flags().StringArrayVar(&values, "value", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	buildCmd.Flags().StringArrayVar(&stringValues, "string-value", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the installed rig template to stdout",
	Long: `Builds the installed rig template. Template values are defined in rig.yaml but can
also be defined in command line arguments. Values defined in arguments supersede
values in rig.yaml.

Examples:
rig build
rig build --value deployment.tag=$(git rev-parse HEAD)
rig build --value host=my-app.${CLUSTER}.example.com

	`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		check(err)

		filePath := path.Join(wd, "rig.yaml")

		tmpl, err := template.NewFromFile(filePath)
		check(err)

		err = tmpl.Sync()
		check(err)

		yaml, err := tmpl.Build(filePath, values, stringValues)
		check(err)

		fmt.Print(yaml)
	},
}
