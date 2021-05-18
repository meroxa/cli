package cased

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type AvailableEndpoint string

const (
	APIEndpoint         AvailableEndpoint = "api"
	WorkflowsEndpoint   AvailableEndpoint = "workflows"
	AuditTrailsEndpoint AvailableEndpoint = "audittrails"

	publishURL = "https://publish.cased.com"
	apiURL     = "https://api.cased.com"
)

var (
	// APIURL is the Cased REST API URL (default: https://api.cased.com)
	APIURL = os.Getenv("CASED_API_URL")
	APIKey = os.Getenv("CASED_API_KEY")

	// PublishURL is the Cased API URL for publishing audit trail events
	// (default: https://publish.cased.com)
	PublishURL = os.Getenv("CASED_PUBLISH_URL")

	// PublishKey is the audit trail API key for publishing audit trail events.
	PublishKey = os.Getenv("CASED_PUBLISH_KEY")

	// WorkflowsKey is the workflows API key for managing and triggering workflows.
	WorkflowsKey = os.Getenv("CASED_WORKFLOWS_KEY")
)

var endpoints Endpoints

// Default HTTP client that can be customized by using `SetHTTPClient`.
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// SetHTTPClient sets the global HTTP client.
func SetHTTPClient(client *http.Client) {
	httpClient = client
}

type Endpoints struct {
	API         Endpoint
	AuditTrails Endpoint
	Workflows   Endpoint

	mu sync.RWMutex
}

type Endpoint interface {
	Call(method, path string, params ParamsContainer, i interface{}) error
}

// SetEndpoint lets you configure the endpoint implementation for a particular
// endpoint. Helpful for mocking in tests.
func SetEndpoint(endpoint AvailableEndpoint, e Endpoint) {
	endpoints.mu.Lock()
	defer endpoints.mu.Unlock()

	switch endpoint {
	case APIEndpoint:
		endpoints.API = e
	case AuditTrailsEndpoint:
		endpoints.AuditTrails = e
	case WorkflowsEndpoint:
		endpoints.Workflows = e
	}
}

// GetEndpoint retrieves an endpoint implementation for a particular endpoint.
func GetEndpoint(endpointType AvailableEndpoint) Endpoint {
	var endpoint Endpoint

	endpoints.mu.RLock()
	switch endpointType {
	case APIEndpoint:
		endpoint = endpoints.API
	case AuditTrailsEndpoint:
		endpoint = endpoints.AuditTrails
	case WorkflowsEndpoint:
		endpoint = endpoints.Workflows
	}
	endpoints.mu.RUnlock()
	if endpoint != nil {
		return endpoint
	}

	endpoint = GetEndpointWithConfig(
		endpointType,
		&EndpointConfig{
			HTTPClient: httpClient,
			URL:        nil,
			APIKey:     nil,
		},
	)

	SetEndpoint(endpointType, endpoint)

	return endpoint
}

// GetEndpointWithConfig configures the specified endpoint type with the default
// configuration.
func GetEndpointWithConfig(endpointType AvailableEndpoint, config *EndpointConfig) Endpoint {
	if config.HTTPClient == nil {
		config.HTTPClient = httpClient
	}

	switch endpointType {
	case APIEndpoint:
		if config.URL == nil {
			if APIURL == "" {
				config.URL = String(apiURL)
			} else {
				config.URL = String(APIURL)
			}
		}
		// Prevents double forward slash when constructing API URLs if
		// https://api.cased.com/ is provided as the base URL compared to
		// https://api.cased.com
		config.URL = String(strings.TrimSuffix(*config.URL, "/"))

		if config.APIKey == nil {
			config.APIKey = String(APIKey)
		}

		return newEndpointImplementation(endpointType, config)

	case AuditTrailsEndpoint:
		if config.URL == nil {
			if PublishURL == "" {
				config.URL = String(publishURL)
			} else {
				config.URL = String(PublishURL)
			}
		}
		config.URL = String(strings.TrimSuffix(*config.URL, "/"))

		if config.APIKey == nil {
			config.APIKey = String(PublishKey)
		}

	case WorkflowsEndpoint:
		if config.URL == nil {
			if APIURL == "" {
				config.URL = String(apiURL)
			} else {
				config.URL = String(APIURL)
			}
		}
		config.URL = String(strings.TrimSuffix(*config.URL, "/"))

		if config.APIKey == nil {
			config.APIKey = String(WorkflowsKey)
		}

		return newEndpointImplementation(endpointType, config)
	}

	return nil
}

type EndpointConfig struct {
	HTTPClient *http.Client
	APIKey     *string
	URL        *string
}

func newEndpointImplementation(endpointType AvailableEndpoint, config *EndpointConfig) Endpoint {
	return &EndpointImplementation{
		HTTPClient: config.HTTPClient,
		Endpoint:   endpointType,
		URL:        *config.URL,
		APIKey:     *config.APIKey,
	}
}

type EndpointImplementation struct {
	Endpoint   AvailableEndpoint
	HTTPClient *http.Client
	URL        string
	APIKey     string
}

func (ei *EndpointImplementation) Call(method, path string, params ParamsContainer, i interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	url := ei.URL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "cased-go/v0.1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ei.APIKey))

	resp, err := ei.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// After a resource is successfully deleted an error or decoding of the
	// response is not necessary
	if method == http.MethodDelete && resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode >= 400 {
		return decodeErrorResponse(resp)
	}

	return json.NewDecoder(resp.Body).Decode(i)
}

func decodeErrorResponse(resp *http.Response) error {
	apiError := &Error{}
	if err := json.NewDecoder(resp.Body).Decode(apiError); err != nil {
		return err
	}

	return apiError
}

type ParamsContainer interface {
	GetParams() *Params
}

type Params struct {
	Context context.Context `json:"-"`
}

func (p *Params) GetParams() *Params {
	return p
}

// Built in type pointer helpers
func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

// Publisher describes the interface for structs that want to publish audit
// events to Cased.
type Publisher interface {
	Publish(event AuditEvent) error
	Options() PublisherOptions
	Flush(timeout time.Duration) bool
}

// Publish publishes an audit event to Cased.
func Publish(event AuditEvent) error {
	client := CurrentPublisher()
	if client.Options().Silence {
		Logger.Println("Audit event was silenced.")
		return nil
	}

	return client.Publish(event)
}

// PublishWithContext enriches the provided audit event with the context set in
// the request. If the same key is present in both the context and provided
// audit event, the audit event value will be preserved.
func PublishWithContext(ctx context.Context, event AuditEvent) error {
	c := GetContextFromContext(ctx)
	for key, value := range c {
		if _, ok := event[key]; ok {
			continue
		}

		event[key] = value
	}

	return Publish(event)
}

// Flush waits for audit events to be published.
func Flush(timeout time.Duration) bool {
	return CurrentPublisher().Flush(timeout)
}
