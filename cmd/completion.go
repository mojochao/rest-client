package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Output completion script",
	Long: strings.Replace(`To load completions:

Bash:

  $ source <(APP_NAME completion bash)
  
To load completions for each session, execute once:
  
Linux:
  $ APP_NAME completion bash > /etc/bash_completion.d/APP_NAME
  
MacOS:
  $ APP_NAME completion bash > /usr/local/etc/bash_completion.d/APP_NAME


Zsh:

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for each session, execute once:

  $ APP_NAME completion zsh > "${fpath[1]}/_APP_NAME"

You will need to start a new shell for this setup to take effect.


Fish:

  $ APP_NAME completion fish | source

To load completions for each session, execute once:

  $ APP_NAME completion fish > ~/.config/fish/completions/APP_NAME.fish

`, "APP_NAME", appName, -1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			return
		}

		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
