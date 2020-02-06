package meroxa

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	apiURL      = "https://api.meroxa.io/v1"
	contentType = "application/json"
)

// Client represents the Meroxa API Client
type Client struct {
	BaseURL   *url.URL
	userAgent string
	username  string
	password  string

	httpClient *http.Client
}

// New returns a configured Meroxa API Client
func New(username, password, ua string) (*Client, error) {
	u, err := url.Parse(getAPIURL())
	if err != nil {
		return nil, err
	}
	c := &Client{
		BaseURL:    u,
		userAgent:  userAgent(ua),
		username:   username,
		password:   password,
		httpClient: httpClient(),
	}
	return c, nil
}

func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	// Merge params
	if params != nil {
		q := req.URL.Query()
		for k, v := range params { // v is a []string
			for _, vv := range v {
				q.Add(k, vv)
			}
			req.URL.RawQuery = q.Encode()
		}
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set Basic Auth
	req.SetBasicAuth(c.username, c.password)

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	req.Header.Add("User-Agent", c.userAgent)
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func getAPIURL() string {
	if u := os.Getenv("API_URL"); u != "" {
		return u
	}

	return apiURL
}

func userAgent(ua string) string {
	if ua != "" {
		return ua
	}
	return "meroxa-go"
}
