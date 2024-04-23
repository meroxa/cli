package global

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	jsonContentType = "application/json"
)

//go:generate mockgen -source=basic_client.go -package=mock -destination=mock/basic_client_mock.go
type BasicClient interface {
	CollectionRequestMultipart(context.Context, string, string, string, interface{}, url.Values, map[string]string) (*http.Response, error)
	CollectionRequest(context.Context, string, string, string, interface{}, url.Values) (*http.Response, error)
	URLRequest(context.Context, string, string, interface{}, url.Values, http.Header) (*http.Response, error)
	AddHeader(string, string)
	SetTimeout(time.Duration)
	ResetBaseURL() error
}

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
	headers    http.Header
	userAgent  string
}

func NewBasicClient() (BasicClient, error) {
	// @TODO incorporate tenant subdomain
	baseURL := GetMeroxaTenantURL()
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

	headers := make(http.Header)
	headers.Add("Meroxa-CLI-Version", Version)

	r := &client{
		baseURL:   u,
		userAgent: fmt.Sprintf("Meroxa CLI %s", Version),
		httpClient: &http.Client{
			Timeout:   flagTimeout,
			Transport: transport,
		},
		headers: headers,
	}
	return r, nil
}

func (c *client) ResetBaseURL() error {
	baseURL := GetMeroxaTenantURL()
	if flagAPIURL != "" {
		baseURL = flagAPIURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	c.baseURL = u
	return nil
}

func (c *client) AddHeader(key, value string) {
	if len(c.headers) == 0 {
		c.headers = make(http.Header)
	}
	c.headers.Add(key, value)
}

func (c *client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

func (c *client) CollectionRequest(
	ctx context.Context,
	method string,
	collection string,
	id string,
	body interface{},
	params url.Values,
) (*http.Response, error) {
	path := fmt.Sprintf("/api/collections/%s/records", collection)
	if len(id) != 0 {
		path += fmt.Sprintf("/%s", id)
	}

	req, err := c.newRequest(ctx, method, path, body, params, nil)
	if err != nil {
		return nil, err
	}

	// Merge params
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) CollectionRequestMultipart(
	ctx context.Context,
	method, collection, id string,
	body interface{},
	params url.Values,
	files map[string]string,
) (*http.Response, error) {
	path := fmt.Sprintf("/api/collections/%s/records", collection)
	if id != "" {
		path += fmt.Sprintf("/%s", id)
	}
	req, err := c.newRequestMultiPart(ctx, method, path, body, params, nil, files)
	if err != nil {
		return nil, err
	}
	// Merge params
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) URLRequest(
	ctx context.Context,
	method, path string,
	body interface{},
	params url.Values,
	headers http.Header,
) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body, params, headers)
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

//nolint:gocyclo
func (c *client) newRequestMultiPart(
	ctx context.Context,
	method string,
	path string,
	body interface{},
	params url.Values,
	headers http.Header,
	files map[string]string,
) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)

	bodyMP := make(map[string]interface{})
	byteData, _ := json.Marshal(body)
	if err = json.Unmarshal(byteData, &bodyMP); err != nil {
		return nil, err
	}

	var w io.Writer
	for k, v := range bodyMP {
		if v != "" {
			w, err = mp.CreateFormField(k)
			if err != nil {
				return nil, err
			}

			if err = c.encodeBody(w, v); err != nil {
				return nil, err
			}
		}
	}

	var file *os.File
	for k, v := range files {
		file, err = os.Open(v)
		if err != nil {
			return nil, err
		}

		w, err = mp.CreateFormFile(k, file.Name())
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(w, file); err != nil {
			return nil, err
		}
	}

	mp.Close()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// add global headers to request
	if len(c.headers) > 0 {
		req.Header = c.headers
	}

	// No need to check for a valid token when trying to authenticate.
	// TODO: Need to change this once we integrate with OAuth2
	if path != "/api/collections/users/auth-with-password" {
		accessToken, err := GetUserToken()
		if err != nil {
			return nil, err
		}
		if _, exists := req.Header["Authorization"]; !exists {
			req.Header.Set("Authorization", accessToken)
		}
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", mp.FormDataContentType())
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

func (c *client) newRequest(
	ctx context.Context,
	method string,
	path string,
	body interface{},
	params url.Values,
	headers http.Header,
) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if encodeErr := c.encodeBody(buf, body); encodeErr != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// add global headers to request
	if len(c.headers) > 0 {
		req.Header = c.headers
	}

	// No need to check for a valid token when trying to authenticate.
	// TODO: Need to change this once we integrate with OAuth2
	if path != "/api/collections/users/auth-with-password" {
		accessToken, err := GetUserToken()
		if err != nil {
			return nil, err
		}
		if _, exists := req.Header["Authorization"]; !exists {
			req.Header.Set("Authorization", accessToken)
		}
	}
	req.Header.Add("Content-Type", jsonContentType)
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("User-Agent", c.userAgent)
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

func (c *client) encodeBody(w io.Writer, v interface{}) error {
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

func GetUserToken() (accessToken string, err error) {
	accessToken = Config.GetString(AccessTokenEnv)
	if accessToken == "" {
		// we need at least one token for creating an authenticated client
		return "", errors.New("please login or signup by running 'meroxa login'")
	}

	return accessToken, nil
}
