package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitbucketApi_CopyPermission(t *testing.T) {
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
						objectId:         "{abcd-1234-5678-90ef}",
						objectName:       "update",
						objectType:       ObjectTypeUser,
						permissionBefore: PermissionTypeRead,
						permissionAfter:  PermissionTypeWrite,
						update:           true,
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
						objectId:         "developer",
						objectName:       "g1",
						objectType:       ObjectTypeGroup,
						permissionBefore: PermissionTypeRead,
						permissionAfter:  PermissionTypeAdmin,
						update:           true,
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
						objectId:         "{1111-aaaa-bbbb-cccc}",
						objectName:       "user3",
						objectType:       ObjectTypeUser,
						permissionBefore: PermissionTypeAdmin,
						remove:           true,
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
			err := ba.CopyPermission(ctx, tt.args.workspace, tt.args.repository, tt.args.operations)
			assert.NoError(t, err)
		})
	}
}
