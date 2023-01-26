package api

import (
	"fmt"
	"sort"
	"strings"
)

type Operation struct {
	objectId         string
	objectName       string
	objectType       ObjectType
	permissionBefore PermissionType
	permissionAfter  PermissionType

	add    bool
	remove bool
	update bool
}

type OperationType string

const (
	OperationTypeAdd    = "add"
	OperationTypeRemove = "remove"
	OperationTypeUpdate = "update"
)

func NewOperationFromPermission(p Permission, operationType OperationType) Operation {
	return Operation{
		objectId:        p.ObjectId,
		objectName:      p.ObjectName,
		objectType:      p.ObjectType,
		permissionAfter: p.Permission,
		add:             operationType == OperationTypeAdd,
		remove:          operationType == OperationTypeRemove,
		update:          operationType == OperationTypeUpdate,
	}
}

func (o Operation) Same() bool {
	return !o.add && !o.update && !o.remove
}

func (o Operation) Message() string {
	switch {
	case o.update:
		return fmt.Sprintf("Update: %s %s %s => %s", o.objectType, o.objectName, strings.ToUpper(string(o.permissionBefore)), strings.ToUpper(string(o.permissionAfter)))
	case o.add:
		return fmt.Sprintf("Add: %s %s (%s)", o.objectType, o.objectName, strings.ToUpper(string(o.permissionAfter)))
	case o.remove:
		return fmt.Sprintf("Remove: %s %s (%s)", o.objectType, o.objectName, strings.ToUpper(string(o.permissionBefore)))
	default:
		return fmt.Sprintf("Same: %s %s (%s)", o.objectType, o.objectName, strings.ToUpper(string(o.permissionAfter)))
	}
}

func MakeOperationList(srcPermissions, targetPermissions []Permission) []Operation {
	srcPermissionsMap := map[string]Permission{}
	for _, v := range srcPermissions {
		srcPermissionsMap[v.ObjectId] = v
	}
	targetPermissionsMap := map[string]Permission{}
	for _, v := range targetPermissions {
		targetPermissionsMap[v.ObjectId] = v
	}

	operations := make(map[string]Operation, 0)

	for k, vs := range srcPermissionsMap {
		if vt, ok := targetPermissionsMap[k]; ok {
			operations[k] = Operation{
				objectId:         vs.ObjectId,
				objectName:       vs.ObjectName,
				objectType:       vs.ObjectType,
				permissionBefore: vt.Permission,
				permissionAfter:  vs.Permission,
				update:           vt.Permission != vs.Permission,
			}
		} else {
			operations[k] = Operation{
				objectId:         vs.ObjectId,
				objectName:       vs.ObjectName,
				objectType:       vs.ObjectType,
				permissionBefore: "",
				permissionAfter:  vs.Permission,
				add:              true,
			}
		}
	}
	for k, vt := range targetPermissionsMap {
		if _, ok := operations[k]; ok {
			continue
		}
		if _, ok := srcPermissionsMap[k]; !ok {
			operations[k] = Operation{
				objectId:         vt.ObjectId,
				objectName:       vt.ObjectName,
				objectType:       vt.ObjectType,
				permissionBefore: vt.Permission,
				permissionAfter:  "",
				remove:           true,
			}
		}
	}

	result := make([]Operation, 0)
	for _, v := range operations {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]
		return strings.Compare(string(a.objectType), string(b.objectType)) < 0 || strings.Compare(string(a.objectId), string(b.objectId)) < 0
	})

	return result
}
