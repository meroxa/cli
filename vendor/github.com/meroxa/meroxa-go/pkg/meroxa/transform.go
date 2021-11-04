package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Property struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

// Transform represent the Meroxa Transform type within the Meroxa API
type Transform struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Required    bool       `json:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Properties  []Property `json:"properties"`
}

// ListTransforms returns an array of Transforms (scoped to the calling user)
func (c *client) ListTransforms(ctx context.Context) ([]*Transform, error) {
	path := fmt.Sprintf("/v1/transforms")

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var tt []*Transform
	err = json.NewDecoder(resp.Body).Decode(&tt)
	if err != nil {
		return nil, err
	}

	return tt, nil
}
