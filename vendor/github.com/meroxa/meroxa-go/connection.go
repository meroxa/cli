package meroxa

import "context"

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
	return nil, nil
}

// GetConnection returns a Connector for the given connection ID
func (c *Client) GetConnection(ctx context.Context, id int) (*Connector, error) {
	return nil, nil
}

// DeleteConnection deletes the Connector with the given id
func (c *Client) DeleteConnection(ctx context.Context, id int) error {
	return nil
}
