package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pipelinesBasePath = "/v1/pipelines"

type PipelineState string

const (
	PipelineStateHealthy  PipelineState = "healthy"
	PipelineStateDegraded PipelineState = "degraded"
)

type PipelineIdentifier struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Pipeline represents the Meroxa Pipeline type within the Meroxa API
type Pipeline struct {
	CreatedAt   time.Time              `json:"created_at"`
	Environment *EntityIdentifier      `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Name        string                 `json:"name"`
	State       PipelineState          `json:"state"`
	UpdatedAt   time.Time              `json:"updated_at"`
	UUID        string                 `json:"uuid"`
}

// CreatePipelineInput represents the input when creating a Meroxa Pipeline
type CreatePipelineInput struct {
	Name        string                 `json:"name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Environment *EntityIdentifier      `json:"environment,omitempty"`
}

// UpdatePipelineInput represents the input when updating a Meroxa Pipeline
type UpdatePipelineInput struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreatePipeline provisions a new Pipeline
func (c *client) CreatePipeline(ctx context.Context, input *CreatePipelineInput) (*Pipeline, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, pipelinesBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdatePipeline updates a pipeline
func (c *client) UpdatePipeline(ctx context.Context, pipelineNameOrID string, input *UpdatePipelineInput) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%s", pipelinesBasePath, pipelineNameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodPatch, path, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// @TODO implement pipeline actions
// UpdatePipelineStatus updates the status of a pipeline
func (c *client) UpdatePipelineStatus(ctx context.Context, pipelineNameOrID string, action Action) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%s/status", pipelinesBasePath, pipelineNameOrID)

	cr := struct {
		State Action `json:"state,omitempty"`
	}{
		State: action,
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, path, cr, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// ListPipelines returns an array of Pipelines (scoped to the calling user)
func (c *client) ListPipelines(ctx context.Context) ([]*Pipeline, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, pipelinesBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var pp []*Pipeline
	err = json.NewDecoder(resp.Body).Decode(&pp)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

// GetPipelineByName returns a Pipeline with the given name
func (c *client) GetPipelineByName(ctx context.Context, name string) (*Pipeline, error) {
	params := map[string][]string{
		"name": {name},
	}

	resp, err := c.MakeRequest(ctx, http.MethodGet, pipelinesBasePath, nil, params)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPipeline returns a Pipeline with the given id
func (c *client) GetPipeline(ctx context.Context, pipelineID int) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%d", pipelinesBasePath, pipelineID)
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// DeletePipeline deletes the Pipeline with the given id
func (c *client) DeletePipeline(ctx context.Context, nameOrID string) error {
	path := fmt.Sprintf("%s/%s", pipelinesBasePath, nameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return err
	}

	return nil
}
