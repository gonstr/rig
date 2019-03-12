package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gonstr/rig/pkg/template"
)

func init() {
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
		tmpl, err := template.NewTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tmpl.Sync()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tmpl.Install()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
