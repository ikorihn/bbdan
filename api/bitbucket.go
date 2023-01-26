package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const urlBitbucketApi = "https://api.bitbucket.org/2.0"

const (
	endpointPermissionConfigUsers  = "/repositories/%s/%s/permissions-config/users"
	endpointPermissionConfigGroups = "/repositories/%s/%s/permissions-config/groups"
	endpointPermissionConfigUser   = "/repositories/%s/%s/permissions-config/users/%s"
	endpointPermissionConfigGroup  = "/repositories/%s/%s/permissions-config/groups/%s"
)

type BitbucketApi struct {
	hc *http.Client

	baseUrl  string
	username string
	password string
}

func NewBitbucketApi(
	hc *http.Client,
	username string,
	password string,
) *BitbucketApi {
	return &BitbucketApi{
		hc:       hc,
		baseUrl:  urlBitbucketApi,
		username: username,
		password: password,
	}
}

type ObjectType string

const (
	ObjectTypeUser  ObjectType = "user"
	ObjectTypeGroup ObjectType = "group"
)

type PermissionType string

const (
	PermissionTypeRead  PermissionType = "read"
	PermissionTypeWrite PermissionType = "write"
	PermissionTypeAdmin PermissionType = "admin"
)

type Permission struct {
	ObjectId   string
	ObjectName string
	ObjectType ObjectType
	Permission PermissionType
}

// APIレスポンス

// errorResponse
type errorResponse struct {
	Type  string     `json:"type"`
	Error errorField `json:"error"`
}
type errorField struct {
	Fields  map[string]string `json:"fields,omitempty"`
	Message string            `json:"message"`
}

type response[T any] struct {
	Values  []T     `json:"values"`
	Pagelen int     `json:"pagelen"`
	Size    int     `json:"size"`
	Next    *string `json:"next,omitempty"`
}

type repositoryPermissionUser struct {
	Type       string        `json:"type"`
	Permission string        `json:"permission"`
	User       bitbucketUser `json:"user"`
}

type repositoryPermissionGroup struct {
	Type       string         `json:"type"`
	Permission string         `json:"permission"`
	Group      bitbucketGroup `json:"group"`
}

type bitbucketUser struct {
	Type        string `json:"type"`
	Uuid        string `json:"uuid"`
	AccountId   string `json:"account_id"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
}

type bitbucketGroup struct {
	Type     string `json:"type"`
	Slug     string `json:"slug"`
	FullSlug string `json:"full_slug"`
	Name     string `json:"name"`
}

func (ba BitbucketApi) do(ctx context.Context, endpoint, method string, body io.Reader) ([]byte, error) {
	u, err := url.Parse(ba.baseUrl + endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(ba.username, ba.password)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := ba.hc.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		var e errorResponse
		err = json.Unmarshal(b, &e)
		if err != nil {
			e = errorResponse{
				Error: errorField{
					Message: string(b),
				},
			}
		}

		return nil, fmt.Errorf("http request error: %v, %v", res.Status, e)
	}

	return b, nil
}

// ListPermission gets permissions for a repository.
func (ba *BitbucketApi) ListPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	res, err := ba.do(ctx, fmt.Sprintf(endpointPermissionConfigGroups, workspace, repository), "GET", nil)
	if err != nil {
		return nil, err
	}
	var groupPermission response[repositoryPermissionGroup]
	err = json.Unmarshal(res, &groupPermission)
	if err != nil {
		return nil, err
	}
	res, err = ba.do(ctx, fmt.Sprintf(endpointPermissionConfigUsers, workspace, repository), "GET", nil)
	if err != nil {
		return nil, err
	}
	var userPermission response[repositoryPermissionUser]
	err = json.Unmarshal(res, &userPermission)
	if err != nil {
		return nil, err
	}

	permissions := make([]Permission, 0)
	for _, v := range groupPermission.Values {
		p := Permission{
			ObjectId:   v.Group.Slug,
			ObjectName: v.Group.Name,
			ObjectType: ObjectType(v.Group.Type),
			Permission: PermissionType(v.Permission),
		}
		permissions = append(permissions, p)
	}
	for _, v := range userPermission.Values {
		p := Permission{
			ObjectId:   v.User.Uuid,
			ObjectName: v.User.Nickname,
			ObjectType: ObjectType(v.User.Type),
			Permission: PermissionType(v.Permission),
		}
		permissions = append(permissions, p)
	}

	return permissions, nil

}

// UpdatePermissions updates permissions of a repository according to operations.
func (ba *BitbucketApi) UpdatePermissions(ctx context.Context, workspace, repository string, operations []Operation) error {
	for _, v := range operations {
		endpoint := ""
		switch {
		case v.update, v.add:
			if v.objectType == ObjectTypeUser {
				endpoint = fmt.Sprintf(endpointPermissionConfigUser, workspace, repository, v.objectId)
			} else {
				endpoint = fmt.Sprintf(endpointPermissionConfigGroup, workspace, repository, v.objectId)
			}

			body, err := json.Marshal(map[string]string{
				"permission": string(v.permissionAfter),
			})
			if err != nil {
				return err
			}
			_, err = ba.do(ctx, endpoint, "PUT", bytes.NewBuffer(body))
			if err != nil {
				return err
			}

		case v.remove:
			if v.objectType == ObjectTypeUser {
				endpoint = fmt.Sprintf(endpointPermissionConfigUser, workspace, repository, v.objectId)
			} else {
				endpoint = fmt.Sprintf(endpointPermissionConfigGroup, workspace, repository, v.objectId)
			}
			_, err := ba.do(ctx, endpoint, "DELETE", nil)
			if err != nil {
				return err
			}

		}
	}
	return nil
}
