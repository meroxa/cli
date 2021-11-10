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
	baseURL         = "https://api.meroxa.io"
	jsonContentType = "application/json"
	textContentType = "text/plain"
)

// client represents the Meroxa API Client
type client struct {
	baseURL   *url.URL
	userAgent string
	token     string

	httpClient *http.Client
}

// Client represents the interface to the Meroxa API
type Client interface {
	CreateConnector(ctx context.Context, input *CreateConnectorInput) (*Connector, error)
	DeleteConnector(ctx context.Context, nameOrID string) error
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*Connector, error)
	GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error)
	ListConnectors(ctx context.Context) ([]*Connector, error)
	UpdateConnector(ctx context.Context, nameOrID string, input *UpdateConnectorInput) (*Connector, error)
	UpdateConnectorStatus(ctx context.Context, nameOrID string, state Action) (*Connector, error)

	CreateEndpoint(ctx context.Context, input *CreateEndpointInput) error
	DeleteEndpoint(ctx context.Context, name string) error
	GetEndpoint(ctx context.Context, name string) (*Endpoint, error)
	ListEndpoints(ctx context.Context) ([]Endpoint, error)

	CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*Environment, error)
	DeleteEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error)
	GetEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error)
	ListEnvironments(ctx context.Context) ([]*Environment, error)

	CreatePipeline(ctx context.Context, input *CreatePipelineInput) (*Pipeline, error)
	DeletePipeline(ctx context.Context, id int) error
	GetPipeline(ctx context.Context, pipelineID int) (*Pipeline, error)
	GetPipelineByName(ctx context.Context, name string) (*Pipeline, error)
	ListPipelines(ctx context.Context) ([]*Pipeline, error)
	ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*Connector, error)
	UpdatePipeline(ctx context.Context, pipelineID int, input *UpdatePipelineInput) (*Pipeline, error)
	UpdatePipelineStatus(ctx context.Context, pipelineID int, action Action) (*Pipeline, error)

	CreateResource(ctx context.Context, input *CreateResourceInput) (*Resource, error)
	DeleteResource(ctx context.Context, nameOrID string) error
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*Resource, error)
	ListResources(ctx context.Context) ([]*Resource, error)
	UpdateResource(ctx context.Context, nameOrID string, input *UpdateResourceInput) (*Resource, error)
	RotateTunnelKeyForResource(ctx context.Context, nameOrID string) (*Resource, error)
	ValidateResource(ctx context.Context, nameOrID string) (*Resource, error)

	ListResourceTypes(ctx context.Context) ([]string, error)

	ListTransforms(ctx context.Context) ([]*Transform, error)

	GetUser(ctx context.Context) (*User, error)

	MakeRequest(ctx context.Context, method string, path string, body interface{}, params url.Values) (*http.Response, error)
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
func New(options ...Option) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &client{
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

func (c *client) MakeRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Response, error) {
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

func (c *client) newRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Request, error) {
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
