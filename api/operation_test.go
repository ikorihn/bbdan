package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperation(t *testing.T) {
	type args struct {
		srcPermissions    []Permission
		targetPermissions []Permission
	}

	tests := []struct {
		name        string
		args        args
		want        []Operation
		wantMessage []string
	}{

		{
			name: "normal",
			args: args{
				srcPermissions: []Permission{
					{
						ObjectId:   "{abc}",
						ObjectName: "same",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeRead,
					},
					{
						ObjectId:   "{abcd-1234-5678-90ef}",
						ObjectName: "update",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeWrite,
					},
					{
						ObjectId:   "developer",
						ObjectName: "diff-grp",
						ObjectType: ObjectTypeGroup,
						PermissionType: PermissionTypeRead,
					},
					{
						ObjectId:   "{xyz}",
						ObjectName: "add",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeWrite,
					},
				},
				targetPermissions: []Permission{
					{
						ObjectId:   "{abcd-1234-5678-90ef}",
						ObjectName: "update",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeRead,
					},
					{
						ObjectId:   "{abc}",
						ObjectName: "same",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeRead,
					},
					{
						ObjectId:   "{lmn}",
						ObjectName: "remove",
						ObjectType: ObjectTypeUser,
						PermissionType: PermissionTypeAdmin,
					},
					{
						ObjectId:   "developer",
						ObjectName: "diff-grp",
						ObjectType: ObjectTypeGroup,
						PermissionType: PermissionTypeAdmin,
					},
				},
			},
			want: []Operation{
				{
					objectId:         "developer",
					objectName:       "diff-grp",
					objectType:       ObjectTypeGroup,
					permissionCurrent: PermissionTypeAdmin,
					permissionAfter:  PermissionTypeRead,
					update:           true,
				},
				{
					objectId:         "{abcd-1234-5678-90ef}",
					objectName:       "update",
					objectType:       ObjectTypeUser,
					permissionCurrent: PermissionTypeRead,
					permissionAfter:  PermissionTypeWrite,
					update:           true,
				},
				{
					objectId:         "{abc}",
					objectName:       "same",
					objectType:       ObjectTypeUser,
					permissionCurrent: PermissionTypeRead,
					permissionAfter:  PermissionTypeRead,
				},
				{
					objectId:         "{lmn}",
					objectName:       "remove",
					objectType:       ObjectTypeUser,
					permissionCurrent: PermissionTypeAdmin,
					permissionAfter:  "",
					remove:           true,
				},
				{
					objectId:         "{xyz}",
					objectName:       "add",
					objectType:       ObjectTypeUser,
					permissionCurrent: "",
					permissionAfter:  PermissionTypeWrite,
					add:              true,
				},
			},
			wantMessage: []string{
				"Update: group diff-grp ADMIN => READ",
				"Update: user update READ => WRITE",
				"Same: user same (READ)",
				"Remove: user remove (ADMIN)",
				"Add: user add (WRITE)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeOperationList(tt.args.srcPermissions, tt.args.targetPermissions)
			assert.Equal(t, tt.want, got)

			fmt.Printf("%+v", got)

			for i, v := range got {
				assert.Equal(t, tt.wantMessage[i], v.Message())
			}
		})
	}
}
