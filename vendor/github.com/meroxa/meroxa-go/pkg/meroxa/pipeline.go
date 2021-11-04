package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const pipelinesBasePath = "/v1/pipelines"

type PipelineState string

const (
	PipelineStateHealthy  PipelineState = "healthy"
	PipelineStateDegraded PipelineState = "degraded"
)

// Pipeline represents the Meroxa Pipeline type within the Meroxa API
type Pipeline struct {
	ID    int           `json:"id"`
	Name  string        `json:"name"`
	State PipelineState `json:"state"`
}

type CreatePipelineInput struct {
	Name string `json:"name"`
}

type UpdatePipelineInput struct {
	Name string `json:"name"`
}

// ComponentKind enum for Component "kinds" within Pipeline stages
type ComponentKind int

const (
	// ConnectorComponent is a Pipeline stage component of type Connector
	ConnectorComponent ComponentKind = 0

	// FunctionComponent is a Pipeline stage component of type Function
	FunctionComponent ComponentKind = 1
)

// PipelineStage represents the Meroxa PipelineStage type within the Meroxa API
type PipelineStage struct {
	ID            int           `json:"id"`
	PipelineID    int           `json:"pipeline_id"`
	ComponentID   int           `json:"component_id"`
	ComponentKind ComponentKind `json:"component_kind"`
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
func (c *client) UpdatePipeline(ctx context.Context, pipelineID int, input *UpdatePipelineInput) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%d", pipelinesBasePath, pipelineID)

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
func (c *client) UpdatePipelineStatus(ctx context.Context, pipelineID int, action Action) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%d/status", pipelinesBasePath, pipelineID)

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

// GetPipelineByName returns a Pipeline with the given name (scoped to the calling user)

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
func (c *client) DeletePipeline(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s/%d", pipelinesBasePath, id)

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
