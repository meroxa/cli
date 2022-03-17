package meroxa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ResourcesBasePath = "/v1/resources"

type ResourceState string

const (
	ResourceStatePending  ResourceState = "pending"
	ResourceStateStarting ResourceState = "starting"
	ResourceStateError    ResourceState = "error"
	ResourceStateReady    ResourceState = "ready"
)

var ErrMissingScheme = errors.New("URL scheme required")

// Credentials represents the Meroxa Resource credentials type within the
// Meroxa API
type Credentials struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	CACert        string `json:"ca_cert"`
	ClientCert    string `json:"client_cert"`
	ClientCertKey string `json:"client_cert_key"`
	UseSSL        bool   `json:"ssl"`
}

// CreateResourceInput represents the input for a Meroxa Resource type we're creating within the Meroxa API
type CreateResourceInput struct {
	Credentials *Credentials            `json:"credentials,omitempty"`
	Environment *EntityIdentifier       `json:"environment,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	Name        string                  `json:"name,omitempty"`
	SSHTunnel   *ResourceSSHTunnelInput `json:"ssh_tunnel,omitempty"`
	Type        ResourceType            `json:"type"`
	URL         string                  `json:"url"`
}

type ResourceSSHTunnelInput struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type ResourceSSHTunnel struct {
	Address   string `json:"address"`
	PublicKey string `json:"public_key"`
}

type ResourceStatus struct {
	State         ResourceState `json:"state"`
	Details       string        `json:"details"`
	LastUpdatedAt time.Time     `json:"last_updated_at"`
}

// Resource represents the Meroxa Resource type within the Meroxa API
type Resource struct {
	UUID        string                 `json:"uuid"`
	Type        ResourceType           `json:"type"`
	Name        string                 `json:"name"`
	URL         string                 `json:"url"`
	Credentials *Credentials           `json:"credentials,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	SSHTunnel   *ResourceSSHTunnel     `json:"ssh_tunnel,omitempty"`
	Environment *EntityIdentifier      `json:"environment,omitempty"`
	Status      ResourceStatus         `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// UpdateResourceInput represents the Meroxa Resource we're updating in the Meroxa API
type UpdateResourceInput struct {
	Name        string                  `json:"name,omitempty"`
	URL         string                  `json:"url,omitempty"`
	Credentials *Credentials            `json:"credentials,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	SSHTunnel   *ResourceSSHTunnelInput `json:"ssh_tunnel,omitempty"`
}

// CreateResource provisions a new Resource from the given CreateResourceInput struct
func (c *client) CreateResource(ctx context.Context, input *CreateResourceInput) (*Resource, error) {
	// url encode url username/password if needed
	var err error
	input.URL, err = encodeURLCreds(input.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, ResourcesBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var r Resource
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *client) UpdateResource(ctx context.Context, nameOrID string, input *UpdateResourceInput) (*Resource, error) {
	// url encode url username/password if needed
	var err error

	if input.URL != "" {
		input.URL, err = encodeURLCreds(input.URL)

		if err != nil {
			return nil, err
		}
	}

	resp, err := c.MakeRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", ResourcesBasePath, nameOrID), input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var r Resource
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *client) RotateTunnelKeyForResource(ctx context.Context, nameOrID string) (*Resource, error) {
	return c.performResourceAction(ctx, nameOrID, "rotate_keys")
}

func (c *client) ValidateResource(ctx context.Context, nameOrID string) (*Resource, error) {
	return c.performResourceAction(ctx, nameOrID, "validate")
}

func (c *client) performResourceAction(ctx context.Context, nameOrID string, action string) (*Resource, error) {
	path := fmt.Sprintf("%s/%s/actions", ResourcesBasePath, nameOrID)
	body := struct {
		Action string `json:"action"`
	}{
		Action: action,
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, path, body, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var rr Resource
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return &rr, nil
}

// ListResources returns an array of Resources (scoped to the calling user)
func (c *client) ListResources(ctx context.Context) ([]*Resource, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, ResourcesBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var rr []*Resource
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

// GetResourceByNameOrID returns a Resource with the given identifier
func (c *client) GetResourceByNameOrID(ctx context.Context, nameOrID string) (*Resource, error) {
	path := fmt.Sprintf("%s/%s", ResourcesBasePath, nameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var r Resource
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// DeleteResource deletes the Resource with the given id
func (c *client) DeleteResource(ctx context.Context, nameOrID string) error {
	path := fmt.Sprintf("%s/%s", ResourcesBasePath, nameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return err
	}

	return nil
}

// Reassemble URL in order to properly encode username and password
func encodeURLCreds(u string) (string, error) {
	s1 := strings.SplitAfter(u, "://")
	scheme := s1[0] // pull out scheme
	if len(s1) == 1 {
		return "", ErrMissingScheme
	}

	v := strings.Split(s1[1], "@") // pull out everything after the @
	if len(v) == 1 {               // no username and password
		return u, nil
	}

	rest := v[len(v)-1]
	userInfoPart := strings.Join(v[:len(v)-1], "@")

	escapedURL, err := url.Parse(scheme + rest)
	if err != nil {
		return "", err
	}

	if rest != "" {
		userinfo := strings.Split(userInfoPart, ":")
		if len(userinfo) > 1 {
			escapedURL.User = url.UserPassword(userinfo[0], userinfo[1])
		} else {
			escapedURL.User = url.UserPassword(userinfo[0], "")
		}
	}

	return escapedURL.String(), nil
}
