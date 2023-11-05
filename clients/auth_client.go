package apiclients

import (
	"context"
	"github.com/a-novel/go-apis"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

type UserTokenHeader struct {
	IAT time.Time `json:"iat"`
	EXP time.Time `json:"exp"`
	ID  uuid.UUID `json:"id"`
}

type UserTokenPayload struct {
	ID uuid.UUID `json:"id"`
}

type UserToken struct {
	Header  UserTokenHeader  `json:"header"`
	Payload UserTokenPayload `json:"payload"`
}

type UserTokenStatus struct {
	OK        bool       `json:"ok"`
	Expired   bool       `json:"expired"`
	NotIssued bool       `json:"notIssued"`
	Malformed bool       `json:"malformed"`
	Token     *UserToken `json:"token,omitempty"`
	TokenRaw  string     `json:"tokenRaw,omitempty"`
}

type AuthClient interface {
	IntrospectToken(ctx context.Context, token string) (*UserTokenStatus, error)
	Ping(ctx context.Context) error
}

type authClientImpl struct {
	url *url.URL
}

func NewAuthClient(url *url.URL) AuthClient {
	return &authClientImpl{url: url}
}

func (a *authClientImpl) IntrospectToken(ctx context.Context, token string) (*UserTokenStatus, error) {
	output := new(UserTokenStatus)
	return output, apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/auth"),
		Method:          http.MethodGet,
		Headers:         map[string]string{"Authorization": token},
		SuccessStatuses: []int{http.StatusOK},
	}, output)
}

func (a *authClientImpl) Ping(ctx context.Context) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/ping"),
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, nil)
}
