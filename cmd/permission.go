package cmd

import (
	"github.com/spf13/cobra"
)

// permissionCmd represents the permission command
var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "List, operate permission of repository",
}

func init() {
	rootCmd.AddCommand(permissionCmd)
}
