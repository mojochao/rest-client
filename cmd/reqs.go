package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reqsCmd represents the reqs command
var reqsCmd = &cobra.Command{
	Use:   "reqs",
	Short: "List rest-client requests",
	Long: `List rest-client requests

If no --http-file options are provided, looks for any files in the current
working directory with the '.http' extension.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load requests from http files.
		allReqs, err := loadReqs(httpFiles)
		if err != nil {
			return err
		}

		// Success! Print requests.
		for i, req := range allReqs {
			fmt.Print(req)
			if i < len(allReqs)-1 {
				fmt.Print("\n###\n\n")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reqsCmd)
}
