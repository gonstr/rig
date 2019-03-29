package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

var templatePath string
var values []string
var stringValues []string

func init() {
	buildCmd.Flags().StringVar(&templatePath, "path", "", "set the template path (this must be a local file path)")
	buildCmd.Flags().StringArrayVar(&values, "value", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	buildCmd.Flags().StringArrayVar(&stringValues, "string-value", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds rig.yaml to stdout",
	Long: `Builds rig.yaml to stdout. Template path and values are defined in
rig.yaml but can also be defined as command line arguments. Values defined in
arguments supersede values in rig.yaml.

If path is specified as a command line argument there is no need for a rig.yaml
file.

Examples:
rig build
rig build --value deployment.tag=$(git rev-parse HEAD)
rig build --path manifests --value host=my-app.${CLUSTER}.example.com

	`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		check(err)

		filePath := path.Join(wd, "rig.yaml")

		tmpl, err := template.New(filePath, templatePath)
		check(err)

		err = tmpl.Sync()
		check(err)

		yaml, err := tmpl.Build(values, stringValues)
		check(err)

		fmt.Print(yaml)
	},
}
