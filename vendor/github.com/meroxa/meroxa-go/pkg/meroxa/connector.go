package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	// ConnectorTypeSource should be changed these since they are associated to `Type` where in fact we're representing
	// source or destination connectors as part of their metadata not this attribute.
	// `Type` should be simply a string (`jdbc-source`, `s3-destination`)
	// They're currently being used in the CLI.
	ConnectorTypeSource      ConnectorType = "source"
	ConnectorTypeDestination ConnectorType = "destination"
)

type Connector struct {
	Configuration map[string]interface{} `json:"config"`
	CreatedAt     time.Time              `json:"created_at"`
	Environment   *EntityIdentifier      `json:"environment,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
	Name          string                 `json:"name"`
	PipelineName  string                 `json:"pipeline_name"`
	ResourceName  string                 `json:"resource_name"`
	Streams       map[string]interface{} `json:"streams"`
	State         ConnectorState         `json:"state"`
	Trace         string                 `json:"trace,omitempty"`
	Type          ConnectorType          `json:"type"`
	UpdatedAt     time.Time              `json:"updated_at"`
	UUID          string                 `json:"uuid"`
}

type CreateConnectorInput struct {
	Name          string                 `json:"name,omitempty"`
	ResourceName  string                 `json:"resource_name"`
	PipelineName  string                 `json:"pipeline_name"`
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
func (c *client) ListPipelineConnectors(ctx context.Context, pipelineNameOrID string) ([]*Connector, error) {
	path := fmt.Sprintf("%s/%s/connectors", pipelinesBasePath, pipelineNameOrID)

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

func filterConnectorsPerType(list []*Connector, cType ConnectorType) []*Connector {
	connectors := make([]*Connector, 0)
	for _, connector := range list {
		if connector.Metadata != nil &&
			connector.Metadata["mx:connectorType"] == string(cType) {
			connectors = append(connectors, connector)
		}
	}
	return connectors
}
