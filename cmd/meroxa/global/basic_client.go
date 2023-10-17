package global

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	jsonContentType = "application/json"
)

type BasicClient interface {
	CollectionRequest(ctx context.Context, method string, collection string, body interface{}, params url.Values, headers http.Header, output interface{}) (*http.Response, error)
	UrlRequest(ctx context.Context, method, path string, body interface{}, params url.Values, headers http.Header, output interface{}) (*http.Response, error)
}

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
	headers    http.Header
	userAgent  string
}

func NewBasicClient() (BasicClient, error) {
	// @TODO incorporate tenant subdomain
	baseURL := GetMeroxaAPIURL()
	if flagAPIURL != "" {
		baseURL = flagAPIURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport
	if flagDebug {
		transport = &dumpTransport{
			out:                    os.Stdout,
			transport:              transport,
			obfuscateAuthorization: true,
		}
	}
	timeout := 5 * time.Second
	if flagTimeout != 0 {
		timeout = flagTimeout
	}
	headers := make(http.Header)
	headers.Add("Meroxa-CLI-Version", Version)

	r := &client{
		baseURL:   u,
		userAgent: fmt.Sprintf("Meroxa CLI %s", Version),
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		headers: headers,
	}
	return r, nil
}

func (r *client) CollectionRequest(ctx context.Context, method, collection string, body interface{}, params url.Values, headers http.Header, output interface{}) (*http.Response, error) {
	path := fmt.Sprintf("/api/collections/%s/records", collection)
	req, err := r.newRequest(ctx, method, path, body, params, headers)
	if err != nil {
		return nil, err
	}

	// Merge params
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	if output != nil {
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (r *client) UrlRequest(ctx context.Context, method, path string, body interface{}, params url.Values, headers http.Header, output interface{}) (*http.Response, error) {
	req, err := r.newRequest(ctx, method, path, body, params, headers)
	if err != nil {
		return nil, err
	}

	// Merge params
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	if output != nil {
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (r *client) newRequest(ctx context.Context, method, path string, body interface{}, params url.Values, headers http.Header) (*http.Request, error) {
	u, err := r.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if err := r.encodeBody(buf, body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// add global headers to request
	if len(r.headers) > 0 {
		req.Header = r.headers
	}
	req.Header.Add("Authorization", getAuthToken())
	req.Header.Add("Content-Type", jsonContentType)
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("User-Agent", r.userAgent)
	for key, value := range headers {
		req.Header.Add(key, strings.Join(value, ","))
	}

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

func (r *client) encodeBody(w io.Writer, v interface{}) error {
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