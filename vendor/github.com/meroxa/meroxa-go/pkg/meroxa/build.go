package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const buildsBasePath = "/v1/builds"

type CreateBuildInput struct {
	SourceBlob SourceBlob `json:"source_blob"`
}

type BuildStatus struct {
	State   string `json:"state"`
	Details string `json:"details"`
}

type Build struct {
	Uuid       string      `json:"uuid"`
	Status     BuildStatus `json:"status"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
	SourceBlob SourceBlob  `json:"source_blob"`
	Image      string      `json:"image"`
}

func (c *client) GetBuild(ctx context.Context, uuid string) (*Build, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", buildsBasePath, uuid), nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var b *Build
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *client) CreateBuild(ctx context.Context, input *CreateBuildInput) (*Build, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, buildsBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var b *Build
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
