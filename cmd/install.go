package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

var force bool

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "FORCE install even if a template has already been installed. This will overwrite rig.yaml")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a remote rig template to the current directory",
	Long: `Installs a rig template from a remote github repository to the current directory.
Template data will be stored in rig.yaml. Git branch/tag or commit can be defined
as a fragment in the template url.

Examples:

rig install https://github.com/gonstr/rig-templates/simple-app
rig install https://github.com/gonstr/rig-templates/simple-app#master
rig install https://github.com/gonstr/rig-templates/simple-app#simple-app/v1.0.0
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires a template uri arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, err := template.NewFromURL(args[0], "", nil)
		check(err)

		err = tmpl.Sync()
		check(err)

		err = tmpl.Install(force)
		check(err)

		fmt.Println("Template installed. Edit values in rig.yaml to your liking and run 'rig build' to build the template")
	},
}
