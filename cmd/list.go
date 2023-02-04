package cmd

import (
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List permissions of a bitbucket repository",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("List permissions for %s/%s\n", workspace, repository)

		hc := http.DefaultClient

		ba := api.NewBitbucketApi(hc, username, password)
		showPermissions(ba, workspace, repository)

		return nil
	},
}

func init() {
	permissionCmd.AddCommand(listCmd)
}
