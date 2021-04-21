package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const endpointBasePath = "/v1/endpoints"

type createEndpointRequest struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Stream   string `json:"stream"`
}

type Endpoint struct {
	Name              string `json:"name"`
	Protocol          string `json:"protocol"`
	Host              string `json:"host"`
	Stream            string `json:"stream"`
	Ready             bool   `json:"ready"`
	BasicAuthUsername string `json:"basic_auth_username"`
	BasicAuthPassword string `json:"basic_auth_password"`
}

func (c *Client) CreateEndpoint(ctx context.Context, name, protocol, stream string) error {
	body := createEndpointRequest{
		Name:     name,
		Protocol: protocol,
		Stream:   stream,
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, endpointBasePath, body, nil)
	if err != nil {
		return err
	}

	return handleAPIErrors(resp)
}

func (c *Client) GetEndpoint(ctx context.Context, name string) (*Endpoint, error) {
	path := fmt.Sprintf("%s/%s", endpointBasePath, name)
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var end Endpoint
	err = json.NewDecoder(resp.Body).Decode(&end)
	if err != nil {
		return nil, err
	}

	return &end, nil
}

func (c *Client) DeleteEndpoint(ctx context.Context, name string) error {
	path := fmt.Sprintf("%s/%s", endpointBasePath, name)
	resp, err := c.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	return handleAPIErrors(resp)
}

func (c *Client) ListEndpoints(ctx context.Context) ([]Endpoint, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, endpointBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var ends []Endpoint
	err = json.NewDecoder(resp.Body).Decode(&ends)
	if err != nil {
		return nil, err
	}

	return ends, nil
}
