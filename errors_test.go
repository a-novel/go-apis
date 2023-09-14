package apis_test

import (
	"fmt"
	"github.com/a-novel/go-apis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorToHTTPCode(t *testing.T) {
	data := []struct {
		name string

		err             error
		errorsToCode    []apis.HTTPError
		catchHTTPErrors bool

		expectStatus int
	}{
		{
			name: "Success",
			err:  fooErr,
			errorsToCode: []apis.HTTPError{
				{
					Err:  barErr,
					Code: http.StatusNotFound,
				},
				{
					Err:  fooErr,
					Code: http.StatusBadRequest,
				},
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Success/NotFound",
			err:  fooErr,
			errorsToCode: []apis.HTTPError{
				{
					Err:  barErr,
					Code: http.StatusNotFound,
				},
			},
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "Success/CatchHTTPErrors",
			err:  &apis.HTTPClientErr{Code: http.StatusUnauthorized, Err: fooErr},
			errorsToCode: []apis.HTTPError{
				{
					Err:  barErr,
					Code: http.StatusNotFound,
				},
			},
			catchHTTPErrors: true,
			expectStatus:    http.StatusUnauthorized,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			apis.ErrorToHTTPCode(c, d.err, d.errorsToCode, d.catchHTTPErrors)
			require.Equal(t, d.expectStatus, w.Code)
		})
	}
}

func TestNewHTTPClientErr(t *testing.T) {
	data := []struct {
		name string

		res            *http.Response
		expectedStatus []int

		expect *apis.HTTPClientErr
	}{
		{
			name: "Success",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400"),
			},
		},
		{
			name: "Success/WithBody",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("response body content")),
			},
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400: %w", fmt.Errorf("response body content")),
			},
		},
		{
			name: "Success/WithExpectedStatus",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expectedStatus: []int{http.StatusOK},
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/WithExpectedStatuses",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expectedStatus: []int{http.StatusOK, http.StatusNotFound},
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected any status in [200 404], got 400"),
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := apis.NewHTTPClientErr(d.res, d.expectedStatus...)
			require.Equal(t, d.expect, err)
		})
	}
}

func TestAsHTTPClientErr(t *testing.T) {
	data := []struct {
		name string

		err error

		expect *apis.HTTPClientErr
	}{
		{
			name: "Success",
			err: apis.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			}, http.StatusOK),
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/WithBody",
			err: apis.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("response body content")),
			}),
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400: %w", fmt.Errorf("response body content")),
			},
		},
		{
			name: "Success/Nested",
			err: fmt.Errorf("error wrapper message: %w", apis.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			}, http.StatusOK)),
			expect: &apis.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/NotAnHTTPClientErr",
			err:  fmt.Errorf("not an HTTPClientErr"),
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := apis.AsHTTPClientErr(d.err)
			require.Equal(t, d.expect, err)
		})
	}
}
