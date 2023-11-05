package apiclients

import (
	"context"
	"github.com/a-novel/go-apis"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

type UpdateImproveRequestVotesForm struct {
	ID        uuid.UUID `json:"id" form:"id"`
	UserID    uuid.UUID `json:"userID" form:"userID"`
	UpVotes   int       `json:"upVotes" form:"upVotes"`
	DownVotes int       `json:"downVotes" form:"downVotes"`
}

type UpdateImproveSuggestionVotesForm struct {
	ID        uuid.UUID `json:"id" form:"id"`
	UserID    uuid.UUID `json:"userID" form:"userID"`
	UpVotes   int       `json:"upVotes" form:"upVotes"`
	DownVotes int       `json:"downVotes" form:"downVotes"`
}

type GetImproveRequestQuery struct {
	ID uuid.UUID `json:"id" form:"id"`
}

type GetImproveSuggestionQuery struct {
	ID uuid.UUID `json:"id" form:"id"`
}

type ImproveRequestPreview struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	UserID  uuid.UUID `json:"userID"`
	Title   string    `json:"title"`
	Content string    `json:"content"`

	UpVotes   int `json:"upVotes"`
	DownVotes int `json:"downVotes"`

	SuggestionsCount         int `json:"suggestionsCount"`
	AcceptedSuggestionsCount int `json:"acceptedSuggestionsCount"`
	RevisionCount            int `json:"revisionsCount"`
}

type ImproveSuggestion struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	SourceID  uuid.UUID `json:"sourceID"`
	UserID    uuid.UUID `json:"userID"`
	Validated bool      `json:"validated"`

	UpVotes   int `json:"upVotes"`
	DownVotes int `json:"downVotes"`

	RequestID uuid.UUID `json:"requestID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type ForumClient interface {
	VoteImproveRequest(ctx context.Context, form UpdateImproveRequestVotesForm) error
	VoteImproveSuggestion(ctx context.Context, form UpdateImproveSuggestionVotesForm) error
	GetImproveRequest(ctx context.Context, query GetImproveRequestQuery) (*ImproveRequestPreview, error)
	GetImproveSuggestion(ctx context.Context, query GetImproveSuggestionQuery) (*ImproveSuggestion, error)

	Ping(ctx context.Context) error
}

type forumClientImpl struct {
	url *url.URL
}

func NewForumClient(url *url.URL) ForumClient {
	return &forumClientImpl{url: url}
}

func (a *forumClientImpl) VoteImproveRequest(ctx context.Context, form UpdateImproveRequestVotesForm) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/improve-request/vote"),
		Method:          http.MethodPost,
		Body:            form,
		SuccessStatuses: []int{http.StatusNoContent},
	}, nil)
}

func (a *forumClientImpl) VoteImproveSuggestion(ctx context.Context, form UpdateImproveSuggestionVotesForm) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/improve-suggestion/vote"),
		Method:          http.MethodPost,
		Body:            form,
		SuccessStatuses: []int{http.StatusNoContent},
	}, nil)
}

func (a *forumClientImpl) GetImproveRequest(ctx context.Context, query GetImproveRequestQuery) (*ImproveRequestPreview, error) {
	pathQuery := url.Values{}
	pathQuery.Set("id", query.ID.String())

	path := a.url.JoinPath("/improve-request")
	path.RawQuery = pathQuery.Encode()

	output := new(ImproveRequestPreview)

	return output, apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            path,
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, output)
}

func (a *forumClientImpl) GetImproveSuggestion(ctx context.Context, query GetImproveSuggestionQuery) (*ImproveSuggestion, error) {
	pathQuery := url.Values{}
	pathQuery.Set("id", query.ID.String())

	path := a.url.JoinPath("/improve-suggestion")
	path.RawQuery = pathQuery.Encode()

	output := new(ImproveSuggestion)

	return output, apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            path,
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, output)
}

func (a *forumClientImpl) Ping(ctx context.Context) error {
	return apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/ping"),
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, nil)
}
