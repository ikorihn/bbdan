package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const urlBitbucketApi = "https://api.bitbucket.org/2.0"

const (
	endpointPermissionConfigUsers  = "/repositories/%s/%s/permissions-config/users"
	endpointPermissionConfigGroups = "/repositories/%s/%s/permissions-config/groups"
	endpointPermissionConfigUser   = "/repositories/%s/%s/permissions-config/users/%s"
	endpointPermissionConfigGroup  = "/repositories/%s/%s/permissions-config/groups/%s"
	endpointDefaultReviewers       = "/repositories/%s/%s/default-reviewers"
	endpointDefaultReviewer        = "/repositories/%s/%s/default-reviewers/%s"
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
	ObjectId       string
	ObjectName     string
	ObjectType     ObjectType
	PermissionType PermissionType
}

type Account struct {
	Uuid        string `json:"uuid"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
}

// Response from Bitbucket API

type errorResponse struct {
	Type  string     `json:"type"`
	Error errorField `json:"error"`
}
type errorField struct {
	Fields  map[string]string `json:"fields,omitempty"`
	Message string            `json:"message"`
}

// response is successful response that has values(list of object)
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

	fmt.Printf("request to %s\n", u.String())

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

// ListGroupPermission gets group permissions for a repository.
func (ba *BitbucketApi) ListGroupPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	permissions := make([]Permission, 0)
	next := fmt.Sprintf(endpointPermissionConfigGroups, workspace, repository)
	for next != "" {
		res, err := ba.do(ctx, next, "GET", nil)
		if err != nil {
			return nil, err
		}
		var groupPermission response[repositoryPermissionGroup]
		err = json.Unmarshal(res, &groupPermission)
		if err != nil {
			return nil, err
		}

		for _, v := range groupPermission.Values {
			p := Permission{
				ObjectId:       v.Group.Slug,
				ObjectName:     v.Group.Name,
				ObjectType:     ObjectType(v.Group.Type),
				PermissionType: PermissionType(v.Permission),
			}
			permissions = append(permissions, p)
		}

		if groupPermission.Next != nil {
			next = *groupPermission.Next
			next = strings.TrimPrefix(next, urlBitbucketApi)
		} else {
			next = ""
		}
	}

	return permissions, nil
}

// ListUserPermission gets group permissions for a repository.
func (ba *BitbucketApi) ListUserPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	permissions := make([]Permission, 0)

	next := fmt.Sprintf(endpointPermissionConfigUsers, workspace, repository)

	for next != "" {
		res, err := ba.do(ctx, next, "GET", nil)
		if err != nil {
			return nil, err
		}
		var userPermission response[repositoryPermissionUser]
		err = json.Unmarshal(res, &userPermission)
		if err != nil {
			return nil, err
		}

		for _, v := range userPermission.Values {
			p := Permission{
				ObjectId:       v.User.Uuid,
				ObjectName:     v.User.Nickname,
				ObjectType:     ObjectType(v.User.Type),
				PermissionType: PermissionType(v.Permission),
			}
			permissions = append(permissions, p)
		}

		if userPermission.Next != nil {
			next = *userPermission.Next
			next = strings.TrimPrefix(next, urlBitbucketApi)
		} else {
			next = ""
		}
	}

	return permissions, nil
}

// ListPermission gets permissions for a repository.
func (ba *BitbucketApi) ListPermission(ctx context.Context, workspace, repository string) ([]Permission, error) {
	permissions := make([]Permission, 0)

	groupPermission, err := ba.ListGroupPermission(ctx, workspace, repository)
	if err != nil {
		return nil, err
	}
	userPermission, err := ba.ListUserPermission(ctx, workspace, repository)
	if err != nil {
		return nil, err
	}
	permissions = append(permissions, groupPermission...)
	permissions = append(permissions, userPermission...)

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

// ListDefaultReviewers gets default reviewers for a repository.
func (ba *BitbucketApi) ListDefaultReviewers(ctx context.Context, workspace, repository string) ([]Account, error) {
	accounts := make([]Account, 0)

	next := fmt.Sprintf(endpointDefaultReviewers, workspace, repository)

	for next != "" {
		res, err := ba.do(ctx, next, "GET", nil)
		if err != nil {
			return nil, err
		}
		var user response[bitbucketUser]
		err = json.Unmarshal(res, &user)
		if err != nil {
			return nil, err
		}

		for _, v := range user.Values {
			a := Account{
				Uuid:        v.Uuid,
				Nickname:    v.Nickname,
				DisplayName: v.DisplayName,
			}
			accounts = append(accounts, a)
		}

		if user.Next != nil {
			next = *user.Next
			next = strings.TrimPrefix(next, urlBitbucketApi)
		} else {
			next = ""
		}
	}

	return accounts, nil
}

// DeleteDefaultReviewers deletes default reviewers for a repository.
// - reviewers: list of the username or the UUID
func (ba *BitbucketApi) DeleteDefaultReviewers(ctx context.Context, workspace, repository string, reviewers []string) ([]Account, error) {
	accounts := make([]Account, 0)

	var wg sync.WaitGroup
	for _, cv := range reviewers {
		wg.Add(1)
		go func(reviewer string) {
			defer wg.Done()

			u := fmt.Sprintf(endpointDefaultReviewer, workspace, repository, reviewer)
			ba.do(ctx, u, "DELETE", nil)
		}(cv)
	}

	wg.Wait()

	return accounts, nil
}

// AddDefaultReviewers adds default reviewers for a repository.
// - reviewers: list of the username or the UUID
func (ba *BitbucketApi) AddDefaultReviewers(ctx context.Context, workspace, repository string, reviewers []string) ([]Account, error) {
	accounts := make([]Account, 0)

	var wg sync.WaitGroup
	for _, cv := range reviewers {
		wg.Add(1)
		go func(reviewer string) {
			defer wg.Done()

			u := fmt.Sprintf(endpointDefaultReviewer, workspace, repository, reviewer)
			ba.do(ctx, u, "PUT", nil)
		}(cv)
	}

	wg.Wait()

	return accounts, nil
}
