package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

var force bool

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "FORCE install even if rig.yaml already exists. This will overwrite rig.yaml")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a rig template in the current directory",
	Long: `Downloads a rig template and installs a .rig file in the
current working directory.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Requires a template uri arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, err := template.NewFromURI(args[0])
		check(err)

		err = tmpl.Sync()
		check(err)

		err = tmpl.Install(force)
		check(err)

		fmt.Println("Template installed. Edit rig.yaml to your liking and run 'rig build' to generate manifests.")
	},
}
