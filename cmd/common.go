package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ikorihn/bbdan/api"
)

func showPermissions(ba *api.BitbucketApi, workspace, repository string) {
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

}

func askOperation(operations []api.Operation) ([]api.Operation, error) {
	messages := make([]string, 0)
	for _, v := range operations {
		if !v.Same() {
			messages = append(messages, v.Message())
		}
	}
	prompt := &survey.MultiSelect{
		Message: "Choose operations:",
		Options: messages,
	}

	selectedIdx := []int{}
	err := survey.AskOne(prompt, &selectedIdx)
	if err != nil {
		return nil, err
	}
	selectedOperations := make([]api.Operation, 0)
	for _, v := range selectedIdx {
		selectedOperations = append(selectedOperations, operations[v])
	}

	return selectedOperations, nil
}
