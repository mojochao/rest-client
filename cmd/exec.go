package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute rest-client requests",
	Long: `Execute rest-client requests.

If no --envs-file options is provided, looks for a 'rest-client.env.json' file
in the current working directory.

If no --http-file options are provided, looks for any files in the current
working directory with the '.http' extension.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load requests from http files.
		allReqs, err := loadReqs(httpFiles)
		if err != nil {
			return err
		}

		// Filter requests by names, if provided.
		if len(reqNames) > 0 {
			allReqs = filterReqs(allReqs, reqNames)
		}

		// Load environment from environments file.
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
		env, exists := envs[envName]
		if !exists {
			return fmt.Errorf("environment not found: %s", envName)
		}

		// Execute requests with environment.
		responses, err := execReqs(allReqs, env)
		if err != nil {
			return err
		}

		// Success! Print responses.
		fmt.Print(renderResponses(responses))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&envName, "env", "e", "", "Environment name")
	execCmd.MarkFlagRequired("env")

	execCmd.Flags().StringArrayVarP(&reqNames, "name", "n", nil, "Request names")
}
