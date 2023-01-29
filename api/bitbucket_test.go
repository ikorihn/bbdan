package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitbucketApi_ListPermissions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repositories/myworkspace/myrepository/permissions-config/groups":
			b, _ := os.ReadFile("testdata/bitbucket_permissions_groups.json")
			var buf bytes.Buffer
			json.Compact(&buf, b)
			res := buf.String()
			p := r.URL.Query().Get("page")
			if p == "" {
				res = strings.TrimSuffix(res, "}") + `,"next": "https://api.bitbucket.org/2.0/repositories/myworkspace/myrepository/permissions-config/groups?page=2" }`
			}
			fmt.Fprint(w, res)

		case "/repositories/myworkspace/myrepository/permissions-config/users":
			b, _ := os.ReadFile("testdata/bitbucket_permissions_users.json")
			var buf bytes.Buffer
			json.Compact(&buf, b)
			res := buf.String()
			p := r.URL.Query().Get("page")
			if p == "" {
				res = strings.TrimSuffix(res, "}") + `,"next": "https://api.bitbucket.org/2.0/repositories/myworkspace/myrepository/permissions-config/users?page=2" }`
			}
			fmt.Fprint(w, res)
		}
	}))
	defer ts.Close()

	ba := &BitbucketApi{
		hc:       http.DefaultClient,
		baseUrl:  ts.URL,
		username: "user",
		password: "pass",
	}
	ctx := context.Background()
	got, err := ba.ListPermission(ctx, "myworkspace", "myrepository")
	want := []Permission{
		{ObjectId: "abcdef1234", ObjectName: "my-group_abcdef1234", ObjectType: "group", PermissionType: "read"},
		{ObjectId: "administrator", ObjectName: "administrator", ObjectType: "group", PermissionType: "admin"},
		{ObjectId: "developer", ObjectName: "developer", ObjectType: "group", PermissionType: "write"},
		{ObjectId: "abcdef1234", ObjectName: "my-group_abcdef1234", ObjectType: "group", PermissionType: "read"},
		{ObjectId: "administrator", ObjectName: "administrator", ObjectType: "group", PermissionType: "admin"},
		{ObjectId: "developer", ObjectName: "developer", ObjectType: "group", PermissionType: "write"},
		{ObjectId: "{1234-fddd-5678-a111}", ObjectName: "john-doe", ObjectType: "user", PermissionType: "admin"},
		{ObjectId: "{aaaa-bbbb-1234-cdef}", ObjectName: "operator-1", ObjectType: "user", PermissionType: "write"},
		{ObjectId: "{9999-9999-9999-9999}", ObjectName: "reader-1", ObjectType: "user", PermissionType: "read"},
		{ObjectId: "{1234-fddd-5678-a111}", ObjectName: "john-doe", ObjectType: "user", PermissionType: "admin"},
		{ObjectId: "{aaaa-bbbb-1234-cdef}", ObjectName: "operator-1", ObjectType: "user", PermissionType: "write"},
		{ObjectId: "{9999-9999-9999-9999}", ObjectName: "reader-1", ObjectType: "user", PermissionType: "read"},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestBitbucketApi_UpdatePermissions(t *testing.T) {
	type args struct {
		workspace  string
		repository string
		operations []Operation
	}

	tests := []struct {
		name       string
		args       args
		wantPath   string
		wantMethod string
		wantBody   []byte
		wantErr    bool
	}{

		{
			name: "update user",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:          "{abcd-1234-5678-90ef}",
						objectName:        "update",
						objectType:        ObjectTypeUser,
						permissionCurrent: PermissionTypeRead,
						permissionAfter:   PermissionTypeWrite,
						update:            true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/users/{abcd-1234-5678-90ef}",
			wantMethod: "PUT",
			wantBody:   []byte(`{"permission":"write"}`),
			wantErr:    false,
		},

		{
			name: "update group",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:          "developer",
						objectName:        "g1",
						objectType:        ObjectTypeGroup,
						permissionCurrent: PermissionTypeRead,
						permissionAfter:   PermissionTypeAdmin,
						update:            true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/groups/developer",
			wantMethod: "PUT",
			wantBody:   []byte(`{"permission":"admin"}`),
			wantErr:    false,
		},

		{
			name: "add user",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:        "{abcd-1234-5678-90ef}",
						objectName:      "user2",
						objectType:      ObjectTypeUser,
						permissionAfter: PermissionTypeAdmin,
						add:             true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/users/{abcd-1234-5678-90ef}",
			wantMethod: "PUT",
			wantBody:   []byte(`{"permission":"admin"}`),
			wantErr:    false,
		},

		{
			name: "add group",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:        "{1234-cdef-5678}",
						objectName:      "g2",
						objectType:      ObjectTypeGroup,
						permissionAfter: PermissionTypeAdmin,
						add:             true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/groups/{1234-cdef-5678}",
			wantMethod: "PUT",
			wantBody:   []byte(`{"permission":"admin"}`),
			wantErr:    false,
		},

		{
			name: "remove user",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:          "{1111-aaaa-bbbb-cccc}",
						objectName:        "user3",
						objectType:        ObjectTypeUser,
						permissionCurrent: PermissionTypeAdmin,
						remove:            true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/users/{1111-aaaa-bbbb-cccc}",
			wantMethod: "DELETE",
			wantBody:   []byte{},
			wantErr:    false,
		},

		{
			name: "remove group",
			args: args{
				workspace:  "myworkspace",
				repository: "myrepository",
				operations: []Operation{
					{
						objectId:        "{9999-9999-9999-9999}",
						objectName:      "g2",
						objectType:      ObjectTypeGroup,
						permissionAfter: PermissionTypeAdmin,
						remove:          true,
					},
				},
			},
			wantPath:   "/repositories/myworkspace/myrepository/permissions-config/groups/{9999-9999-9999-9999}",
			wantMethod: "DELETE",
			wantBody:   []byte{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.wantPath, r.URL.Path)
				assert.Equal(t, tt.wantMethod, r.Method)
				b, _ := io.ReadAll(r.Body)
				r.Body.Close()
				assert.Equal(t, tt.wantBody, b)
			}))
			defer ts.Close()

			ba := &BitbucketApi{
				hc:       http.DefaultClient,
				baseUrl:  ts.URL,
				username: "user",
				password: "pass",
			}
			ctx := context.Background()
			err := ba.UpdatePermissions(ctx, tt.args.workspace, tt.args.repository, tt.args.operations)
			assert.NoError(t, err)
		})
	}
}
