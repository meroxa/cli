package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

// Resource represents the Meroxa Resource type within the Meroxa API
type Resource struct {
	ID            int               `json:"id"`
	Kind          string            `json:"kind"`
	Name          string            `json:"name"`
	URL           string            `json:"url"`
	Credentials   *Credentials      `json:"credentials,omitempty"`
	Configuration map[string]string `json:"configuration"`
	Metadata      map[string]string `json:"metadata"`
}

// CreateResource provisions a new Resource from the given Resource struct
func (c *Client) CreateResource(ctx context.Context, resource *Resource) (*Resource, error) {
	path := fmt.Sprintf("/v1/resources")

	resp, err := c.makeRequest(ctx, http.MethodPost, path, resource, nil)
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
	path := fmt.Sprintf("/v1/resources")

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
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
	path := fmt.Sprintf("/v1/resources/%d", id)

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
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
	path := fmt.Sprintf("/v1/resources")

	params := map[string][]string{
		"name": []string{name},
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, params)
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
	path := fmt.Sprintf("/v1/resources/%d", id)

	_, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// ListResourceConnections returns an array of Connectors for a given Resource
func (c *Client) ListResourceConnections(ctx context.Context, id int) ([]*Connector, error) {
	path := fmt.Sprintf("/v1/resources/%d/connections", id)

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
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
