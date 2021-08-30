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
	ID          int                     `json:"id"`
	Type        string                  `json:"type"`
	Name        string                  `json:"name,omitempty"`
	URL         string                  `json:"url"`
	Credentials *Credentials            `json:"credentials,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	SSHTunnel   *ResourceSSHTunnelInput `json:"ssh_tunnel,omitempty"`
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
	State         string    `json:"state"`
	Details       string    `json:"details"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// Resource represents the Meroxa Resource type within the Meroxa API
type Resource struct {
	ID          int                    `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	URL         string                 `json:"url"`
	Credentials *Credentials           `json:"credentials,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	SSHTunnel   *ResourceSSHTunnel     `json:"ssh_tunnel,omitempty"`
	Status      ResourceStatus         `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// UpdateResourceInput represents the Meroxa Resource we're updating in the Meroxa API
type UpdateResourceInput struct {
	Name        string                  `json:"name,omitempty"`
	URL         string                  `json:"url,omitempty"`
	Metadata    map[string]interface{}  `json:"metadata,omitempty"`
	Credentials *Credentials            `json:"credentials,omitempty"`
	SSHTunnel   *ResourceSSHTunnelInput `json:"ssh_tunnel,omitempty"`
}

// CreateResource provisions a new Resource from the given CreateResourceInput struct
func (c *Client) CreateResource(ctx context.Context, resource *CreateResourceInput) (*Resource, error) {
	// url encode url username/password if needed
	var err error
	resource.URL, err = encodeURLCreds(resource.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, ResourcesBasePath, resource, nil)
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

func (c *Client) UpdateResource(ctx context.Context, key string, resourceToUpdate UpdateResourceInput) (*Resource, error) {
	// url encode url username/password if needed
	var err error

	if resourceToUpdate.URL != "" {
		resourceToUpdate.URL, err = encodeURLCreds(resourceToUpdate.URL)

		if err != nil {
			return nil, err
		}
	}

	resp, err := c.MakeRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", ResourcesBasePath, key), resourceToUpdate, nil)
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

func (c *Client) RotateTunnelKeyForResource(ctx context.Context, id string) (*Resource, error) {
	return c.performResourceAction(ctx, id, "rotate_keys")
}

func (c *Client) ValidateResource(ctx context.Context, id string) (*Resource, error) {
	return c.performResourceAction(ctx, id, "validate")
}

func (c *Client) performResourceAction(ctx context.Context, id string, action string) (*Resource, error) {
	path := fmt.Sprintf("%s/%s/actions", ResourcesBasePath, id)
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
func (c *Client) ListResources(ctx context.Context) ([]*Resource, error) {
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

// GetResource returns a Resource with the given id
func (c *Client) GetResource(ctx context.Context, id int) (*Resource, error) {
	path := fmt.Sprintf("%s/%d", ResourcesBasePath, id)

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

// GetResourceByName returns a Resource with the given name
func (c *Client) GetResourceByName(ctx context.Context, name string) (*Resource, error) {
	params := map[string][]string{
		"name": []string{name},
	}

	resp, err := c.MakeRequest(ctx, http.MethodGet, ResourcesBasePath, nil, params)
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
func (c *Client) DeleteResource(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s/%d", ResourcesBasePath, id)

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

// ListResourceConnections returns an array of Connectors for a given Resource
func (c *Client) ListResourceConnections(ctx context.Context, id int) ([]*Connector, error) {
	path := fmt.Sprintf("%s/%d/connections", ResourcesBasePath, id)

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var cc []*Connector
	err = json.NewDecoder(resp.Body).Decode(&cc)
	if err != nil {
		return nil, err
	}

	return cc, nil
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
