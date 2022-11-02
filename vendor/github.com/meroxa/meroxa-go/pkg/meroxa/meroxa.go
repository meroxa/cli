package meroxa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL         = "https://api.meroxa.io"
	jsonContentType = "application/json"
	textContentType = "text/plain"
)

// EnvironmentIdentifier represents either value as a unique key for a Meroxa Environment
type EnvironmentIdentifier struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

// EntityIdentifier represents one or both values for a Meroxa Entity
type EntityIdentifier struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

func (e EntityIdentifier) GetNameOrUUID() (string, error) {
	if e.Name != "" {
		return e.Name, nil
	} else if e.UUID != "" {
		return e.UUID, nil
	}
	return "", fmt.Errorf("identifier has neither name or UUID")
}

// client represents the Meroxa API Client
type client struct {
	requester
}

type Requester struct {
	baseURL    *url.URL
	httpClient *http.Client
	headers    http.Header
	userAgent  string
}

type requester interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}, params url.Values, headers http.Header) (*http.Response, error)
}

type account interface {
	ListAccounts(ctx context.Context) ([]*Account, error)
}

// Client represents the interface to the Meroxa API
type Client interface {
	requester
	account

	CreateApplication(ctx context.Context, input *CreateApplicationInput) (*Application, error)
	CreateApplicationV2(ctx context.Context, input *CreateApplicationInput) (*Application, error)
	DeleteApplication(ctx context.Context, name string) error
	DeleteApplicationEntities(ctx context.Context, name string) (*http.Response, error)
	GetApplication(ctx context.Context, name string) (*Application, error)
	ListApplications(ctx context.Context) ([]*Application, error)

	CreateBuild(ctx context.Context, input *CreateBuildInput) (*Build, error)
	GetBuild(ctx context.Context, uuid string) (*Build, error)
	GetBuildLogs(ctx context.Context, uuid string) (*http.Response, error)

	CreateConnector(ctx context.Context, input *CreateConnectorInput) (*Connector, error)
	DeleteConnector(ctx context.Context, nameOrID string) error
	GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*Connector, error)
	GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error)
	ListConnectors(ctx context.Context) ([]*Connector, error)
	UpdateConnector(ctx context.Context, nameOrID string, input *UpdateConnectorInput) (*Connector, error)
	UpdateConnectorStatus(ctx context.Context, nameOrID string, state Action) (*Connector, error)

	GetDeployment(ctx context.Context, appIdentifier string, depUUID string) (*Deployment, error)
	GetLatestDeployment(ctx context.Context, appIdentifier string) (*Deployment, error)
	CreateDeployment(ctx context.Context, input *CreateDeploymentInput) (*Deployment, error)

	CreateFunction(ctx context.Context, input *CreateFunctionInput) (*Function, error)
	GetFunction(ctx context.Context, nameOrUUID string) (*Function, error)
	GetFunctionLogs(ctx context.Context, nameOrUUID string) (*http.Response, error)
	ListFunctions(ctx context.Context) ([]*Function, error)
	DeleteFunction(ctx context.Context, nameOrUUID string) (*Function, error)

	CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*Environment, error)
	DeleteEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error)
	GetEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error)
	UpdateEnvironment(ctx context.Context, nameOrUUID string, input *UpdateEnvironmentInput) (*Environment, error)
	ListEnvironments(ctx context.Context) ([]*Environment, error)
	PerformActionOnEnvironment(ctx context.Context, nameOrUUID string, input *RepairEnvironmentInput) (*Environment, error)

	CreatePipeline(ctx context.Context, input *CreatePipelineInput) (*Pipeline, error)
	DeletePipeline(ctx context.Context, nameOrID string) error
	GetPipeline(ctx context.Context, pipelineID int) (*Pipeline, error)
	GetPipelineByName(ctx context.Context, name string) (*Pipeline, error)
	ListPipelines(ctx context.Context) ([]*Pipeline, error)
	ListPipelineConnectors(ctx context.Context, pipelineNameOrID string) ([]*Connector, error)
	UpdatePipeline(ctx context.Context, pipelineNameOrID string, input *UpdatePipelineInput) (*Pipeline, error)
	UpdatePipelineStatus(ctx context.Context, pipelineNameOrID string, action Action) (*Pipeline, error)

	CreateResource(ctx context.Context, input *CreateResourceInput) (*Resource, error)
	DeleteResource(ctx context.Context, nameOrID string) error
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*Resource, error)
	ListResources(ctx context.Context) ([]*Resource, error)
	UpdateResource(ctx context.Context, nameOrID string, input *UpdateResourceInput) (*Resource, error)
	RotateTunnelKeyForResource(ctx context.Context, nameOrID string) (*Resource, error)
	ValidateResource(ctx context.Context, nameOrID string) (*Resource, error)
	IntrospectResource(ctx context.Context, nameOrID string) (*ResourceIntrospection, error)

	ListResourceTypes(ctx context.Context) ([]string, error)

	CreateSource(ctx context.Context) (*Source, error)

	ListTransforms(ctx context.Context) ([]*Transform, error)

	GetUser(ctx context.Context) (*User, error)
}

// New returns a Meroxa API client. To configure it provide a list of Options.
// Note that by default the client is not using any authentication, to provide
// it please use option WithAuthentication or provide your own *http.Client,
// which takes care of authentication.
//
// Example creating an authenticated client:
//
//	c, err := New(
//	    WithAuthentication(auth.DefaultConfig(), accessToken, refreshToken),
//	)
func New(options ...Option) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	r := &Requester{
		baseURL:   u,
		userAgent: "meroxa-go",
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: http.DefaultTransport,
		},
	}
	for _, opt := range options {
		err := opt(r)
		if err != nil {
			return nil, err
		}
	}
	c := &client{
		requester: r,
	}
	return c, nil
}

func (r *Requester) MakeRequest(ctx context.Context, method, path string, body interface{}, params url.Values, headers http.Header) (*http.Response, error) {
	req, err := r.newRequest(ctx, method, path, body, params, headers)
	if err != nil {
		return nil, err
	}

	// Merge params
	resp, err := r.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Requester) newRequest(ctx context.Context, method, path string, body interface{}, params url.Values, headers http.Header) (*http.Request, error) {
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
	if r.headers != nil {
		req.Header = r.headers
	}
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

func (r *Requester) encodeBody(w io.Writer, v interface{}) error {
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
