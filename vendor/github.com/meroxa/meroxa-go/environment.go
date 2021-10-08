package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

const environmentsBasePath = "/v1/environments"

type EnvironmentStatus struct {
	Details string `json:"details,omitempty"`
	State   string `json:"state"`
}

// Environment represents the Meroxa Environment type within the Meroxa API
type Environment struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Provider string            `json:"provider"`
	Region   string            `json:"region"`
	Status   EnvironmentStatus `json:"status"`
	ID       string            `json:"id"`
}

// ListEnvironments returns an array of Pipelines (scoped to the calling user)
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
