package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const connectorsBasePath = "/v1/connectors"

type ConnectorState string

const (
	ConnectorStatePending ConnectorState = "pending"
	ConnectorStateRunning ConnectorState = "running"
	ConnectorStatePaused  ConnectorState = "paused"
	ConnectorStateCrashed ConnectorState = "crashed"
	ConnectorStateFailed  ConnectorState = "failed"
	ConnectorStateDOA     ConnectorState = "doa"
)

type Action string

const (
	ActionPause   Action = "pause"
	ActionResume  Action = "resume"
	ActionRestart Action = "restart"
)

type ConnectorType string

const (
	ConnectorTypeSource      ConnectorType = "source"
	ConnectorTypeDestination ConnectorType = "destination"
)

type Connector struct {
	ID            int                    `json:"id"`
	Type          ConnectorType          `json:"type"`
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"config"`
	Metadata      map[string]interface{} `json:"metadata"`
	Streams       map[string]interface{} `json:"streams"`
	State         ConnectorState         `json:"state"`
	Trace         string                 `json:"trace,omitempty"`
	PipelineID    int                    `json:"pipeline_id"`
	PipelineName  string                 `json:"pipeline_name"`
}

type CreateConnectorInput struct {
	Name          string                 `json:"name,omitempty"`
	ResourceID    int                    `json:"resource_id"`
	PipelineID    int                    `json:"pipeline_id,omitempty"`
	PipelineName  string                 `json:"pipeline_name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Type          ConnectorType          `json:"connector_type,omitempty"`
	Input         string                 `json:"input,omitempty"`
}

type UpdateConnectorInput struct {
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config,omitempty"`
}

// CreateConnector provisions a connector between the Resource and the Meroxa
// platform
func (c *client) CreateConnector(ctx context.Context, input *CreateConnectorInput) (*Connector, error) {
	if input.Configuration != nil {
		input.Configuration["input"] = input.Input
	} else {
		input.Configuration = map[string]interface{}{"input": input.Input}
	}
	input.Metadata = map[string]interface{}{"mx:connectorType": string(input.Type)}
	resp, err := c.MakeRequest(ctx, http.MethodPost, connectorsBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// @TODO implement connector actions
// UpdateConnectorStatus updates the status of a connector
func (c *client) UpdateConnectorStatus(ctx context.Context, nameOrID string, state Action) (*Connector, error) {
	path := fmt.Sprintf("%s/%s/status", connectorsBasePath, nameOrID)

	cr := struct {
		State Action `json:"state,omitempty"`
	}{
		State: state,
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, path, cr, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// UpdateConnector updates the name, or a configuration of a connector
func (c *client) UpdateConnector(ctx context.Context, nameOrID string, input *UpdateConnectorInput) (*Connector, error) {
	path := fmt.Sprintf("%s/%s", connectorsBasePath, nameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodPatch, path, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// ListPipelineConnectors returns an array of Connectors (scoped to the calling user)
func (c *client) ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*Connector, error) {
	path := fmt.Sprintf("/v1/pipelines/%d/connectors", pipelineID)

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var rr []*Connector
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

// ListConnectors returns an array of Connectors (scoped to the calling user)
func (c *client) ListConnectors(ctx context.Context) ([]*Connector, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, connectorsBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var rr []*Connector
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

// GetConnectorByNameOrID returns a Connector with the given identifier
func (c *client) GetConnectorByNameOrID(ctx context.Context, nameOrID string) (*Connector, error) {
	path := fmt.Sprintf("%s/%s", connectorsBasePath, nameOrID)

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var con Connector
	err = json.NewDecoder(resp.Body).Decode(&con)
	if err != nil {
		return nil, err
	}

	return &con, nil
}

// DeleteConnector deletes the Connector with the given id
func (c *client) DeleteConnector(ctx context.Context, nameOrID string) error {
	path := fmt.Sprintf("%s/%s", connectorsBasePath, nameOrID)

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
