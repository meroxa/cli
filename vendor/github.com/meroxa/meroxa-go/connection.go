package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Connector struct {
	ID            int                    `json:"id"`
	Kind          string                 `json:"type"`
	Name          string                 `json:"name"`
	Configuration map[string]string      `json:"config"`
	Metadata      map[string]string      `json:"metadata"`
	Streams       map[string]interface{} `json:"streams"`
}

// CreateConnection provisions a connection between the Resource and the Meroxa
// platform
func (c *Client) CreateConnection(ctx context.Context, resourceID int, config map[string]string) (*Connector, error) {
	path := fmt.Sprintf("/v1/resources/%d/connection", resourceID)

	options := map[string]map[string]string{
		"configuration": config,
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, path, options, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 204 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Status %d, %v", resp.StatusCode, string(body))
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// ListConnections returns an array of Connections (scoped to the calling user)
func (c *Client) ListConnections(ctx context.Context) ([]*Connector, error) {
	path := fmt.Sprintf("/v1/connectors")

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var rr []*Connector
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

// GetConnection returns a Connector for the given connection ID
func (c *Client) GetConnection(ctx context.Context, id int) (*Connector, error) {
	path := fmt.Sprintf("/v1/connectors/%d", id)

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// GetConnectionByName returns a Connection with the given name
func (c *Client) GetConnectionByName(ctx context.Context, name string) (*Connector, error) {
	path := fmt.Sprintf("/v1/connectors")

	params := map[string][]string{
		"name": []string{name},
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, params)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// DeleteConnection deletes the Connector with the given id
func (c *Client) DeleteConnection(ctx context.Context, id int) error {
	path := fmt.Sprintf("/v1/connectors/%d", id)

	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("Status %d, %v", resp.StatusCode, err)
	}

	return nil
}
