package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// envsCmd represents the envs command
var envsCmd = &cobra.Command{
	Use:   "envs",
	Short: "List rest-client environments",
	Long: `List rest-client environments defined in environments file.

If no --envs-file options is provided, looks for a 'rest-client.env.json' file
in the current working directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := homedir.Expand(envsFile)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		envs, err := parseEnvs(f)
		if err != nil {
			return err
		}

		for env, vars := range envs {
			fmt.Printf("%s:\n", env)
			for k, v := range vars {
				fmt.Printf("  %s = %s\n", k, v)
			}
			fmt.Print("\n")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envsCmd)
}
