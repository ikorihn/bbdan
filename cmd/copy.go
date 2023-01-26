/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy permissions",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		srcRepository := args[1]
		targetRepository := args[2]
		fmt.Printf("Copy permissions from %s/%s to %s/%s\n", workspace, srcRepository, workspace, targetRepository)

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

		srcPermissions, err := ba.ListPermission(ctx, workspace, srcRepository)
		if err != nil {
			return err
		}
		targetPermissions, err := ba.ListPermission(ctx, workspace, targetRepository)
		if err != nil {
			return err
		}

		operations := api.MakeOperationList(srcPermissions, targetPermissions)
		selectedOperations, err := askOperation(operations)
		if err != nil {
			return err
		}

		err = ba.UpdatePermissions(ctx, workspace, targetRepository, selectedOperations)
		if err != nil {
			fmt.Printf("Failed to update: %v\n", err)
			return err
		}

		showPermissions(ba, workspace, targetRepository)
		return nil
	},
}

func init() {
	permissionCmd.AddCommand(copyCmd)
}
