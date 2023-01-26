package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove permissions",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("Remove selected permissions from %s/%s\n", workspace, repository)

		hc := http.DefaultClient

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
		ba := api.NewBitbucketApi(hc, username, password)
		ctx := context.Background()
		permissions, err := ba.ListPermission(ctx, workspace, repository)
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}

		operations := []api.Operation{}
		for _, v := range permissions {
			o := api.NewOperationFromPermission(v, api.OperationTypeRemove)
			operations = append(operations, o)
		}

		selectedOperations, err := askOperation(operations)
		if err != nil {
			return err
		}
		err = ba.UpdatePermissions(ctx, workspace, repository, selectedOperations)
		if err != nil {
			fmt.Printf("Failed to update: %v\n", err)
			return err
		}

		showPermissions(ba, workspace, repository)

		return nil
	},
}

func init() {
	permissionCmd.AddCommand(removeCmd)
}
