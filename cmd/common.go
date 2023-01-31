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
			v.PermissionType,
		)
	}

}

func askOperation(operations []api.Operation) ([]api.Operation, error) {
	notSame := make([]api.Operation, 0)
	for _, v := range operations {
		if !v.Same() {
			notSame = append(notSame, v)
		}
	}

	messages := make([]string, 0)
	for _, v := range notSame {
		messages = append(messages, v.Message())
	}
	selectedIdx, err := multiSelect("Choose operations:", messages)
	if err != nil {
		return nil, err
	}

	selectedOperations := make([]api.Operation, 0)
	for _, v := range selectedIdx {
		selectedOperations = append(selectedOperations, notSame[v])
	}

	return selectedOperations, nil
}

func askPermissionToUpdate(permissions []api.Permission) ([]api.Permission, error) {
	messages := make([]string, len(permissions))
	for i, v := range permissions {
		messages[i] = fmt.Sprintf("%s %s: %s", v.ObjectType, v.ObjectName, v.PermissionType)
	}
	selectedIdx, err := multiSelect("Choose permissions to update:", messages)
	if err != nil {
		return nil, err
	}

	selected := make([]api.Permission, 0)
	for _, v := range selectedIdx {
		selected = append(selected, permissions[v])
	}

	return selected, nil
}

func askOperationType() (api.OperationType, error) {
	messages := make([]string, 0)
	messages = append(messages, string(api.OperationTypeRemove))
	messages = append(messages, string(api.OperationTypeUpdate))

	prompt := &survey.Select{
		Message: "Choose operation:",
		Options: messages,
	}

	var selected string
	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return "", err
	}

	return api.OperationType(selected), nil
}

func askPermissionType() (api.PermissionType, error) {
	messages := make([]string, 0)
	messages = append(messages, string(api.PermissionTypeAdmin))
	messages = append(messages, string(api.PermissionTypeRead))
	messages = append(messages, string(api.PermissionTypeWrite))

	prompt := &survey.Select{
		Message: "Choose permission:",
		Options: messages,
	}

	var selected string
	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return "", err
	}

	return api.PermissionType(selected), nil
}

func multiSelect(message string, options []string) ([]int, error) {
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	selectedIdx := []int{}
	err := survey.AskOne(prompt, &selectedIdx)
	if err != nil {
		return nil, err
	}

	return selectedIdx, nil
}
