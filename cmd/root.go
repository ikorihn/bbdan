package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bbdan",
	Short: "Unofficial command line tool for Bitbucket Cloud",
}

var (
	username string
	password string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	cobra.OnInitialize(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("$XDG_CONFIG_HOME/bbdan")
		viper.AddConfigPath("$HOME/.config/bbdan")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		username = viper.GetString("username")
		password = viper.GetString("password")
	})

	return rootCmd.Execute()
}
