package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const buildsBasePathV1 = "/v1/builds"
const buildsBasePathV2 = "/v2/builds"

type CreateBuildInput struct {
	SourceBlob  SourceBlob        `json:"source_blob"`
	Environment *EntityIdentifier `json:"environment,omitempty"`
}

type BuildStatus struct {
	State   string `json:"state"`
	Details string `json:"details"`
}

type Build struct {
	Uuid        string            `json:"uuid"`
	Status      BuildStatus       `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	SourceBlob  SourceBlob        `json:"source_blob"`
	Image       string            `json:"image"`
	Environment *EntityIdentifier `json:"environment,omitempty"`
}

func (c *client) GetBuild(ctx context.Context, uuid string) (*Build, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", buildsBasePathV1, uuid), nil, nil, nil)
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
	resp, err := c.MakeRequest(ctx, http.MethodPost, buildsBasePathV1, input, nil, nil)
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

func (c *client) GetBuildLogsV2(ctx context.Context, uuid string) (*Logs, error) {
	path := fmt.Sprintf("%s/%s/logs", buildsBasePathV2, uuid)
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var l *Logs
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return nil, err
	}

	return l, nil
}
