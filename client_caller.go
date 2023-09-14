package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"net/http"
	"net/url"
)

type CallHTTPConfig struct {
	// Path is the URL to call.
	Path *url.URL
	// Method is the HTTP method to use.
	Method string
	// Body is the (optional) body of the request. Any valid Go serializable object can be passed.
	Body interface{}
	// Headers are extra headers to send with the request.
	Headers map[string]string
	// SuccessStatuses indicates which status codes should be considered as successful.
	SuccessStatuses []int
	// Client is the http client to use. By default, http.DefaultClient is used.
	Client *http.Client
}

// CallHTTP performs a http call with the given configuration. It automatically parses the result into the output
// pointer, if given.
//
// If the call is successful, but the response code is not valid, a HTTPClientErr is returned.
func CallHTTP(ctx context.Context, cfg CallHTTPConfig, output interface{}) error {
	if cfg.Client == nil {
		cfg.Client = http.DefaultClient
	}

	var reqBody io.Reader

	if cfg.Body != nil {
		bodyBytes, err := json.Marshal(cfg.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.Path.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to initiate request: %w", err)
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	res, err := cfg.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	defer res.Body.Close()

	if _, ok := lo.Find(cfg.SuccessStatuses, func(item int) bool {
		return res.StatusCode == item
	}); !ok {
		return NewHTTPClientErr(res, cfg.SuccessStatuses...)
	}

	if output == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, output); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return nil
}
