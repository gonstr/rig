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

var fromStdin bool
var values []string
var stringValues []string

func init() {
	buildCmd.Flags().BoolVar(&fromStdin, "from-stdin", false, "build template from stdin")
	buildCmd.Flags().StringArrayVar(&values, "value", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	buildCmd.Flags().StringArrayVar(&stringValues, "string-value", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Build a template to stdout",
	Args:  cobra.MaximumNArgs(1),
	Long: `Build a template to stdout.

Template path can be supplied as the first argument. If no argument is supplied,
a rig.yaml file is expected. The template can also be passed to stdin if the
--from-stdin argument is supplied.

Template values can be defined in rig.yaml or by --value or --string-value
arguments. Values supplied in arguments supercede values in rig.yaml.

Example usage:

rig build
rig build --value deployment.tag=$(git rev-parse HEAD)
rig build my/manifests/folder --value host=my-app.${CLUSTER}.example.com
cat manifest.yaml | rig build --from-stdin --string-value port=8080

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if fromStdin {
			bytes, err := ioutil.ReadAll(os.Stdin)
			check(err)

			output, err := build.FromString(string(bytes), nil, values, stringValues)
			check(err)

			fmt.Println(output)
		} else {
			if len(args) > 0 {
				output, err := build.FromTemplatesPath(args[0], nil, values, stringValues)
				check(err)

				fmt.Println(output)
			} else {
				wd, err := os.Getwd()
				check(err)

				rigPath := path.Join(wd, "rig.yaml")

				if !fs.PathExists(rigPath) {
					check(errors.New("invalid command: either supply a template path argument or run the command in a dir with a rig.yaml file"))
				}

				output, err := build.FromRigFile(rigPath, values, stringValues)
				check(err)

				fmt.Println(output)
			}
		}
	},
}
