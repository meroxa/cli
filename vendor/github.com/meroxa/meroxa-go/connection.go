package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Connector struct {
	ID            int               `json:"id"`
	Kind          string            `json:"kind"`
	Name          string            `json:"name"`
	Configuration map[string]string `json:"configuration"`
	Metadata      map[string]string `json:"metadata"`
}

// CreateConnection provisions a connection between the Resource and the Meroxa
// platform
func (c *Client) CreateConnection(ctx context.Context, resourceID int, config map[string]string) (*Connector, error) {
	path := fmt.Sprintf("/v1/resources/%d/connection", resourceID)

	resp, err := c.makeRequest(ctx, http.MethodPost, path, config, nil)
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

// GetConnection returns a Connector for the given connection ID
func (c *Client) GetConnection(ctx context.Context, id int) (*Connector, error) {
	path := fmt.Sprintf("/v1/connections/%d", id)

	resp, err := c.makeRequest(ctx, http.MethodPost, path, nil, nil)
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
	path := fmt.Sprintf("/v1/connections/%d", id)

	_, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
