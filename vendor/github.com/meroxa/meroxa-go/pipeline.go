package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const pipelinesBasePath = "/v1/pipelines"

// Pipeline represents the Meroxa Pipeline type within the Meroxa API
type Pipeline struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata,omitempty"`
	State    string            `json:"state"`
}

type UpdatePipelineInput struct {
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata"`
	State    string            `json:"state"`
}

// ComponentKind enum for Component "kinds" within Pipeline stages
type ComponentKind int

const (
	// ConnectorComponent is a Pipeline stage component of type Connector
	ConnectorComponent ComponentKind = 0

	// FunctionComponent is a Pipeline stage component of type Function
	FunctionComponent = 1
)

// PipelineStage represents the Meroxa PipelineStage type within the Meroxa API
type PipelineStage struct {
	ID            int `json:"id"`
	PipelineID    int `json:"pipeline_id"`
	ComponentID   int `json:"component_id"`
	ComponentKind int `json:"component_kind"`
}

// CreatePipeline provisions a new Pipeline
func (c *Client) CreatePipeline(ctx context.Context, pipeline *Pipeline) (*Pipeline, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, pipelinesBasePath, pipeline, nil)
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
func (c *Client) UpdatePipeline(ctx context.Context, pipelineID int, pipelineToUpdate UpdatePipelineInput) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%d", pipelinesBasePath, pipelineID)

	resp, err := c.makeRequest(ctx, http.MethodPatch, path, pipelineToUpdate, nil)
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

// UpdatePipelineStatus updates the status of a pipeline
func (c *Client) UpdatePipelineStatus(ctx context.Context, pipelineID int, state string) (*Pipeline, error) {
	path := fmt.Sprintf("%s/%d/status", pipelinesBasePath, pipelineID)

	cr := struct {
		State string `json:"state,omitempty"`
	}{
		State: state,
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, path, cr, nil)
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
func (c *Client) ListPipelines(ctx context.Context) ([]*Pipeline, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, pipelinesBasePath, nil, nil)
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

// GetPipelineByName returns a Pipeline with the given id
func (c *Client) GetPipelineByName(ctx context.Context, name string) (*Pipeline, error) {
	params := map[string][]string{
		"name": {name},
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, pipelinesBasePath, nil, params)
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

// DeletePipeline deletes the Pipeline with the given id
func (c *Client) DeletePipeline(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s/%d", pipelinesBasePath, id)

	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return err
	}

	return nil
}

// GetPipelineStages returns an array of Pipeline Stages for the given Pipeline id
func (c *Client) GetPipelineStages(ctx context.Context, pipelineID int) ([]*PipelineStage, error) {
	path := fmt.Sprintf("%s/%d/stages", pipelinesBasePath, pipelineID)

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var pp []*PipelineStage
	err = json.NewDecoder(resp.Body).Decode(&pp)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

// AddPipelineStage adds the PipelineStage to the given Pipeline
func (c *Client) AddPipelineStage(ctx context.Context, pipelineID int, connectorID int) (*PipelineStage, error) {
	path := fmt.Sprintf("%s/%d/stages", pipelinesBasePath, pipelineID)

	params := map[string][]string{
		"connector": {strconv.Itoa(connectorID)},
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, path, nil, params)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var s PipelineStage
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// RemovePipelineStage removes the PipelineStage from the given Pipeline
func (c *Client) RemovePipelineStage(ctx context.Context, pipelineID int, stageID int) error {
	path := fmt.Sprintf("%s/%d/stages/%d", pipelinesBasePath, pipelineID, stageID)

	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return err
	}

	return nil
}
