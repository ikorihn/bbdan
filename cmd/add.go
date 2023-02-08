package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add workspace repository user|group id read|write|admin",
	Short: "Add permissions",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("Add permissions to %s/%s\n", workspace, repository)

		objectType := args[2]
		objectId := args[3]
		permissionType := args[4]
		p := api.Permission{
			ObjectType:     api.ObjectType(objectType),
			ObjectId:       objectId,
			PermissionType: api.PermissionType(permissionType),
		}

		hc := http.DefaultClient

		ba := api.NewBitbucketApi(hc, username, password)
		ctx := context.Background()

		operations := []api.Operation{}
		o := api.NewAddOperation(p)
		operations = append(operations, o)

		err := ba.UpdatePermissions(ctx, workspace, repository, operations)
		if err != nil {
			fmt.Printf("Failed to update: %v\n", err)
			return err
		}

		showPermissions(ba, workspace, repository)

		return nil
	},
}

func init() {
	permissionCmd.AddCommand(addCmd)
}
