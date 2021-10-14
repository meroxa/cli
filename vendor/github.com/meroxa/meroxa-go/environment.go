package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const environmentsBasePath = "/v1/environments"

type EnvironmentStatus struct {
	State   string `json:"state"`
	Details string `json:"details,omitempty"`
}

// Environment represents the Meroxa Environment type within the Meroxa API
type Environment struct {
	UUID          string                 `json:"uuid"`
	Name          string                 `json:"name"`
	Provider      string                 `json:"provider"`
	Region        string                 `json:"region"`
	Type          string                 `json:"type"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	Status        EnvironmentStatus      `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// ListEnvironments returns an array of Environments (scoped to the calling user)
func (c *Client) ListEnvironments(ctx context.Context) ([]*Environment, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, environmentsBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var ee []*Environment
	err = json.NewDecoder(resp.Body).Decode(&ee)
	if err != nil {
		return nil, err
	}

	return ee, nil
}
