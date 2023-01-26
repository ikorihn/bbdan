package cmd

import (
	"context"
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
	Run: func(cmd *cobra.Command, args []string) {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("List permissions for %s/%s\n", workspace, repository)

		hc := http.DefaultClient

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		ba := api.NewBitbucketApi(hc, username, password)
		permissions, err := ba.ListPermission(context.Background(), workspace, repository)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		fmt.Println("==== RESULT ====")
		fmt.Println("type, id, name, permission")
		for _, v := range permissions {
			fmt.Printf("%s, %s, %s, %s\n",
				v.ObjectType,
				v.ObjectId,
				v.ObjectName,
				v.Permission,
			)
		}
	},
}

func init() {
	permissionCmd.AddCommand(listCmd)
}
