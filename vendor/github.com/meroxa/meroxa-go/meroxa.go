package meroxa

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL         = "https://api.meroxa.io/v1"
	jsonContentType = "application/json"
	textContentType = "text/plain"
)

// Client represents the Meroxa API Client
type Client struct {
	baseURL   *url.URL
	userAgent string
	token     string

	httpClient *http.Client
}

// New returns a Meroxa API client. To configure it provide a list of Options.
// Note that by default the client is not using any authentication, to provide
// it please use option WithAuthentication or provide your own *http.Client,
// which takes care of authentication.
//
// Example creating an authenticated client:
//  c, err := New(
//      WithAuthentication(auth.DefaultConfig(), accessToken, refreshToken),
//  )
func New(options ...Option) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:   u,
		userAgent: "meroxa-go",
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: http.DefaultTransport,
		},
	}

	for _, opt := range options {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) MakeRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body, params)
	if err != nil {
		return nil, err
	}

	// Merge params
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if err := c.encodeBody(buf, body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// add global headers to request
	req.Header.Add("Content-Type", jsonContentType)
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("User-Agent", c.userAgent)

	// add params
	if params != nil {
		q := req.URL.Query()
		for k, v := range params { // v is a []string
			for _, vv := range v {
				q.Add(k, vv)
			}
			req.URL.RawQuery = q.Encode()
		}
	}

	return req, nil
}

func (c *Client) encodeBody(w io.Writer, v interface{}) error {
	if v == nil {
		return nil
	}

	switch body := v.(type) {
	case string:
		_, err := w.Write([]byte(body))
		return err
	case []byte:
		_, err := w.Write(body)
		return err
	default:
		return json.NewEncoder(w).Encode(v)
	}
}
