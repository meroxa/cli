package meroxa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	ID          int               `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name,omitempty"`
	URL         string            `json:"url"`
	Credentials *Credentials      `json:"credentials,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Resource represents the Meroxa Resource type within the Meroxa API
type Resource struct {
	ID          int               `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Credentials *Credentials      `json:"credentials,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateResourceInput represents the Meroxa Resource we're updating in the Meroxa API
type UpdateResourceInput struct {
	Name        string            `json:"name,omitempty"` // TODO: Update this via CLI
	URL         string            `json:"url,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Credentials *Credentials      `json:"credentials,omitempty"`
}

// CreateResource provisions a new Resource from the given CreateResourceInput struct
func (c *Client) CreateResource(ctx context.Context, resource *CreateResourceInput) (*Resource, error) {
	// url encode url username/password if needed
	var err error
	resource.URL, err = encodeURLCreds(resource.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, ResourcesBasePath, resource, nil)
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

	resp, err := c.makeRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", ResourcesBasePath, key), resourceToUpdate, nil)
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

// ListResources returns an array of Resources (scoped to the calling user)
func (c *Client) ListResources(ctx context.Context) ([]*Resource, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, ResourcesBasePath, nil, nil)
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

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
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

	resp, err := c.makeRequest(ctx, http.MethodGet, ResourcesBasePath, nil, params)
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

	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
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

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
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

	rest := strings.Split(s1[1], "@") // pull out everything after the @
	if len(rest) == 1 {               // no username and password
		return u, nil
	}

	escapedURL, err := url.Parse(scheme + rest[1])
	if err != nil {
		return "", err
	}

	if rest[0] != "" {
		username := strings.Split(rest[0], ":")[0]
		password := strings.Split(rest[0], ":")[1]
		ui := url.UserPassword(username, password)
		escapedURL.User = ui
	}

	return escapedURL.String(), nil
}
