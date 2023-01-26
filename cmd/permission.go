/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// permissionCmd represents the permission command
var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("permission called")
	},
}

func init() {
	rootCmd.AddCommand(permissionCmd)

	permissionCmd.PersistentFlags().StringP("username", "u", "YOUR NAME", "")
	permissionCmd.PersistentFlags().StringP("password", "p", "YOUR PASSWORD", "")

}
