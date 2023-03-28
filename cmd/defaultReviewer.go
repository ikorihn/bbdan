package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ikorihn/bbdan/api"
	"github.com/spf13/cobra"
)

// defaultReviewerCmd represents the defaultReviewer command
var defaultReviewerCmd = &cobra.Command{
	Use:   "default-reviewer",
	Short: "List, operate defaultReviewer of repository",
}

// listDefaultReviewerCmd represents the list command
var listDefaultReviewerCmd = &cobra.Command{
	Use:   "list",
	Short: "List default reviewers of a bitbucket repository",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]
		fmt.Printf("List permissions for %s/%s\n", workspace, repository)

		hc := http.DefaultClient

		ba := api.NewBitbucketApi(hc, username, password)
		ctx := context.Background()
		accounts, err := ba.ListDefaultReviewers(ctx, workspace, repository)
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}

		fmt.Println("==== RESULT ====")
		fmt.Println("id, name")
		for _, v := range accounts {
			fmt.Printf("%s, %s\n",
				v.Uuid,
				v.Nickname,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(defaultReviewerCmd)
	defaultReviewerCmd.AddCommand(listDefaultReviewerCmd)
}
