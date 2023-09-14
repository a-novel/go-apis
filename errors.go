package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type HTTPError struct {
	Err  error
	Code int
}

// ErrorToHTTPCode aborts with a specific status code, depending on the error received. If the error does not match
// any in errorsToCode, http.StatusInternalServerError is returned.
// The error is also automatically added to gin.Context errors.
//
// If catchHTTPErrors is set to true, and the error implements HTTPClientErr, then the code of the HTTPClientErr is
// returned.
func ErrorToHTTPCode(c *gin.Context, err error, errorsToCode []HTTPError, catchHTTPErrors bool) {
	if catchHTTPErrors {
		if httpErr := AsHTTPClientErr(err); httpErr != nil {
			_ = c.AbortWithError(httpErr.Code, httpErr.Err)
			return
		}
	}

	for _, e := range errorsToCode {
		if errors.Is(err, e.Err) {
			_ = c.AbortWithError(e.Code, err)
			return
		}
	}

	_ = c.AbortWithError(http.StatusInternalServerError, err)
}

// HTTPClientErr wraps the error returned by a successful http.Client call, when the response code is not ok.
type HTTPClientErr struct {
	// Err wraps the response body as an error, if any.
	Code int `json:"code"`
	// Code is the response status code.
	Err error `json:"body"`
}

func (h HTTPClientErr) Error() string {
	return fmt.Sprintf("request failed with status %d: %s", h.Code, h.Err)
}

// NewHTTPClientErr wraps the response of a http.Client call into an error object. If available, the response body is
// embedded into the error message.
func NewHTTPClientErr(res *http.Response, expect ...int) error {
	var (
		body       error
		errMessage error
	)

	if res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			body = fmt.Errorf("failed to read response body: %w", err)
		} else {
			body = fmt.Errorf(string(bodyBytes))
		}
	}

	if len(expect) == 0 {
		errMessage = fmt.Errorf("got unexpected status code %d", res.StatusCode)
	} else if len(expect) == 1 {
		errMessage = fmt.Errorf("expected status %d, got %d", expect[0], res.StatusCode)
	} else {
		errMessage = fmt.Errorf("expected any status in %v, got %d", expect, res.StatusCode)
	}

	if body != nil {
		errMessage = fmt.Errorf("%s: %w", errMessage, body)
	}

	return &HTTPClientErr{
		Code: res.StatusCode,
		Err:  errMessage,
	}
}

// AsHTTPClientErr checks if an error implements HTTPClientErr, and return the typed error if so. If err is not
// a HTTPClientErr, nil is returned.
func AsHTTPClientErr(err error) *HTTPClientErr {
	httpClientErr := new(HTTPClientErr)

	if errors.As(err, &httpClientErr) {
		return httpClientErr
	}

	return nil
}
