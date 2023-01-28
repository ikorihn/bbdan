package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update permissions",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("Update selected permissions of %s/%s\n", workspace, repository)

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

		selectedPermissions, err := askPermissionToUpdate(permissions)
		if err != nil {
			return err
		}
		operation, err := askOperationType()
		if err != nil {
			return err
		}

		var permission api.PermissionType
		if operation == api.OperationTypeAdd || operation == api.OperationTypeUpdate {
			permission, err = askPermissionType()
			if err != nil {
				return err
			}
		}

		operations := []api.Operation{}
		for _, v := range selectedPermissions {
			var o api.Operation
			switch operation {
			case api.OperationTypeUpdate:
				o = api.NewUpdateOperation(v, permission)
			case api.OperationTypeRemove:
				o = api.NewRemoveOperation(v)
			default:
				continue
			}
			operations = append(operations, o)
		}

		fmt.Printf("%+v\n", operations)

		// err = ba.UpdatePermissions(ctx, workspace, repository, operations)
		// if err != nil {
		// 	fmt.Printf("Failed to update: %v\n", err)
		// 	return err
		// }

		// showPermissions(ba, workspace, repository)

		return nil
	},
}

func init() {
	permissionCmd.AddCommand(updateCmd)
}
