package api

import (
	"context"
	"encoding/json"
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
	baseUrl string,
	username string,
	password string,
) *BitbucketApi {
	if baseUrl == "" {
		baseUrl = urlBitbucketApi
	}
	return &BitbucketApi{
		hc:       hc,
		baseUrl:  baseUrl,
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
	Type                    string `json:"type"`
	Owner                   string `json:"owner"`
	Workspace               string `json:"workspace"`
	Slug                    string `json:"slug"`
	FullSlug                string `json:"full_slug"`
	Name                    string `json:"name"`
	DefaultPermission       string `json:"default_permission"`
	EmailForwardingDisabled string `json:"email_forwarding_disabled"`
	AccountPrivilege        string `json:"account_privilege"`
}

func (ba *BitbucketApi) ListPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	endpoint, err := url.Parse(ba.baseUrl)
	if err != nil {
		return nil, err
	}
	endpoint = endpoint.JoinPath("repositories", workspace, repository, "permissions-config")

	endpointPermissionGroup := endpoint.JoinPath("groups")
	var groupPermission values[repositoryPermissionGroup]
	err = ba.do(ctx, endpointPermissionGroup.String(), groupPermission)
	if err != nil {
		return nil, err
	}

	endpointPermissionUser := endpoint.JoinPath("users")
	var userPermission values[repositoryPermissionUser]
	err = ba.do(ctx, endpointPermissionUser.String(), userPermission)
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

func (ba BitbucketApi) do(ctx context.Context, url string, v any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := ba.hc.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	return nil
}
