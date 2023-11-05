package apiclients

import (
	"context"
	"github.com/a-novel/go-apis"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type SetUserPermissionsForm struct {
	UserID      uuid.UUID `json:"userID" form:"userID"`
	SetFields   []string  `json:"setFields" form:"setFields"`
	UnsetFields []string  `json:"unsetFields" form:"unsetFields"`
}

type HasUserScopeQuery struct {
	UserID uuid.UUID `json:"userID" form:"userID"`
	Scope  string    `json:"scope" form:"scope"`
}

type PermissionsClient interface {
	SetUserPermissions(ctx context.Context, form SetUserPermissionsForm) error
	HasUserScope(ctx context.Context, query HasUserScopeQuery) error
	Ping(ctx context.Context) error
}

type permissionsClientImpl struct {
	url *url.URL
}

const (
	FieldValidatedAccount = "validated_account"
	FieldAdminAccess      = "admin_access"
)

const (
	CanVotePost              = "forum:post:vote"
	CanPostImproveRequest    = "forum:improve-request:post"
	CanPostImproveSuggestion = "forum:improve-suggestion:post"

	CanUseOpenAIPlayground = "openai:playground"
)

func NewPermissionsClient(url *url.URL) PermissionsClient {
	return &permissionsClientImpl{url: url}
}

func (a *permissionsClientImpl) SetUserPermissions(ctx context.Context, form SetUserPermissionsForm) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/user/permissions"),
		Method:          http.MethodPost,
		Body:            form,
		SuccessStatuses: []int{http.StatusCreated},
	}, nil)
}

func (a *permissionsClientImpl) HasUserScope(ctx context.Context, query HasUserScopeQuery) error {
	pathQuery := new(url.Values)
	pathQuery.Set("userID", query.UserID.String())
	pathQuery.Set("scope", query.Scope)

	path := a.url.JoinPath("/user/scopes")
	path.RawQuery = pathQuery.Encode()

	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            path,
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusNoContent},
	}, nil)
}

func (a *permissionsClientImpl) Ping(ctx context.Context) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/ping"),
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, nil)
}
