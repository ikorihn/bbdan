package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const urlBitbucketApi = "https://api.bitbucket.org/2.0"

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

type values[T any] struct {
	Values []T `json:"values"`
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

func (ba BitbucketApi) do(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(ba.username, ba.password)

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
		type httpErr struct {
			Error map[string]string `json:"error"`
		}
		var e httpErr
		err = json.Unmarshal(b, &e)
		if err != nil {
			e = httpErr{
				Error: map[string]string{
					"message": "unknown",
				},
			}
		}

		return nil, fmt.Errorf("http request error: %v, %v", res.Status, e)
	}

	return b, nil
}

func (ba *BitbucketApi) ListPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	endpoint, err := url.Parse(ba.baseUrl)
	if err != nil {
		return nil, err
	}
	endpoint = endpoint.JoinPath("repositories", workspace, repository, "permissions-config")

	endpointPermissionGroup := endpoint.JoinPath("groups")
	res, err := ba.do(ctx, endpointPermissionGroup.String())
	if err != nil {
		return nil, err
	}
	var groupPermission values[repositoryPermissionGroup]
	err = json.Unmarshal(res, &groupPermission)
	if err != nil {
		return nil, err
	}

	endpointPermissionUser := endpoint.JoinPath("users")
	res, err = ba.do(ctx, endpointPermissionUser.String())
	if err != nil {
		return nil, err
	}
	var userPermission values[repositoryPermissionUser]
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
