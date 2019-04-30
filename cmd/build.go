package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gonstr/rig/pkg/fs"

	"github.com/gonstr/rig/pkg/build"
	"github.com/spf13/cobra"
)

var values []string
var stringValues []string

func init() {
	buildCmd.Flags().StringArrayVar(&values, "value", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	buildCmd.Flags().StringArrayVar(&stringValues, "string-value", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds a template to stdout",
	Args:  cobra.MaximumNArgs(1),
	Long: `Builds a template to stdout.

To build a rig template, rig needs to know where to look for the template files
and what template values to use. There are a few ways to supply this data:

- Template path and values can be defined in rig.yaml as well as in cmd line
  arguments
- Template data can be passed to stdin

Path and values supplied through command line arguments supercede those defined
in rig.yaml.

Example usage:

rig build
rig build --value deployment.tag=$(git rev-parse HEAD)
rig build my/manifests/folder --value host=my-app.${CLUSTER}.example.com
cat manifest.yaml | rig build --string-value port=8080

	`,
	Run: func(cmd *cobra.Command, args []string) {
		fi, err := os.Stdin.Stat()
		check(err)

		if (fi.Mode() & os.ModeCharDevice) == 0 {
			// Data from stdin
			if len(args) > 0 {
				check(errors.New("invalid command: template data passed to stdin AND template path defined as argument"))
			}

			bytes, err := ioutil.ReadAll(os.Stdin)
			check(err)

			output, err := build.FromString(string(bytes), nil, values, stringValues)
			check(err)

			fmt.Println(output)
		} else {
			if len(args) > 0 {
				output, err := build.FromTemplatesPath(args[0], values, stringValues)
				check(err)

				fmt.Println(output)
			} else {
				wd, err := os.Getwd()
				check(err)

				rigPath := path.Join(wd, "rig.yaml")

				if !fs.PathExists(rigPath) {
					check(errors.New("invalid command: either pass a template to stdin, supply a template path argument or run the command in a dir with a rig.yaml file"))
				}

				output, err := build.FromRigFile(rigPath, values, stringValues)
				check(err)

				fmt.Println(output)
			}
		}
	},
}
