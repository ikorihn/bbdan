package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		fmt.Printf("List default reviewers for %s/%s\n", workspace, repository)

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

var overwriteDefaultReviewerCmd = &cobra.Command{
	Use:   "overwrite",
	Short: "Overwrite default reviewer",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace := args[0]
		repository := args[1]

		reviewers := strings.Split(args[2], ",")

		fmt.Printf("Overwrite default reviewers of %s/%s\n", workspace, repository)

		hc := http.DefaultClient

		ba := api.NewBitbucketApi(hc, username, password)
		ctx := context.Background()
		currentReviewers, err := ba.ListDefaultReviewers(ctx, workspace, repository)
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
		curReviewerIds := make([]string, 0)
		for _, v := range currentReviewers {
			curReviewerIds = append(curReviewerIds, v.Uuid)
		}

		ba.DeleteDefaultReviewers(ctx, workspace, repository, curReviewerIds)
		ba.AddDefaultReviewers(ctx, workspace, repository, reviewers)

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
	defaultReviewerCmd.AddCommand(overwriteDefaultReviewerCmd)
}
