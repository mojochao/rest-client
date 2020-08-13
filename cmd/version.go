package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mojochao/rest-client/identity"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long: `Display application version`,
	Run: func(cmd *cobra.Command, args []string) {
		info := identity.GetMetadata()
		extra := ""
		if verbose {
			extra = fmt.Sprintf("(date=%s, branch=%s, commit=%s, state=%s)",
				info.BuildDate, info.GitBranch, info.GitCommit, info.GitState)
		}
		fmt.Printf("%s v%s %s\n", appName, info.Version, extra)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
