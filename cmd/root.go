package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rest-client",
	Short: "CLI REST client compatible with JetBrains REST client",
	Long: `The rest-client application is a command-line interface application that
is capable of executing JetBrains REST client *.http or *.rest request
files against multiple environments defined in rest-client.env.json 
environment files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rest-client.yaml)")
	rootCmd.PersistentFlags().StringVar(&envsFile, "envs-file", defaultEnvFile, "rest-client environments file")
	rootCmd.PersistentFlags().StringArrayVar(&httpFiles, "http-file", httpFiles, "rest-client http files")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rest-client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rest-client")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
