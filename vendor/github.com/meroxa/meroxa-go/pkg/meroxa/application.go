package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApplicationState string

const (
	ApplicationStateInitialized ApplicationState = "initialized"
	ApplicationStateDeploying   ApplicationState = "deploying"
	ApplicationStatePending     ApplicationState = "pending"
	ApplicationStateRunning     ApplicationState = "running"
	ApplicationStateDegraded    ApplicationState = "degraded"
	ApplicationStateFailed      ApplicationState = "failed"
)

const applicationsBasePathV1 = "/v1/applications"
const applicationsBasePathV2 = "/v2/applications"

type ResourceCollection struct {
	Name        string `json:"name,omitempty"`
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
}

type ApplicationResource struct {
	EntityIdentifier
	ResourceType string             `json:"type,omitempty"`
	Status       string             `json:"status,omitempty"`
	Collection   ResourceCollection `json:"collection,omitempty"`
}

type EntityDetails struct {
	EntityIdentifier
	ResourceUUID string `json:"resource_uuid,omitempty"`
	ResourceType string `json:"type,omitempty"`
	Status       string `json:"status,omitempty"`
}

// Application represents the Meroxa Application type within the Meroxa API
type Application struct {
	UUID        string                `json:"uuid"`
	Name        string                `json:"name"`
	Language    string                `json:"language"`
	GitSha      string                `json:"git_sha,omitempty"`
	Status      ApplicationStatus     `json:"status,omitempty"`
	Environment *EntityIdentifier     `json:"environment,omitempty"`
	Pipeline    EntityDetails         `json:"pipeline,omitempty"`
	Connectors  []EntityDetails       `json:"connectors,omitempty"`
	Functions   []EntityDetails       `json:"functions,omitempty"`
	Resources   []ApplicationResource `json:"resources,omitempty"`
	Deployments []EntityIdentifier    `json:"deployments,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	DeletedAt   time.Time             `json:"deleted_at,omitempty"`
}

type ApplicationLogs struct {
	FunctionLogs   map[string]string `json:"functions"`
	ConnectorLogs  map[string]string `json:"connectors"`
	DeploymentLogs map[string]string `json:"latest_deployment"`
}

// CreateApplicationInput represents the input for a Meroxa Application create operation in the API
type CreateApplicationInput struct {
	Name        string            `json:"name"`
	Language    string            `json:"language"`
	GitSha      string            `json:"git_sha,omitempty"`
	Pipeline    EntityIdentifier  `json:"pipeline,omitempty"`
	Environment *EntityIdentifier `json:"environment,omitempty"`
}

type ApplicationStatus struct {
	State   ApplicationState `json:"state"`
	Details string           `json:"details,omitempty"`
}

func (c *client) CreateApplication(ctx context.Context, input *CreateApplicationInput) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, applicationsBasePathV1, input, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var a *Application
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *client) CreateApplicationV2(ctx context.Context, input *CreateApplicationInput) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, applicationsBasePathV2, input, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var a *Application
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, err
	}

	return a, nil
}

func (c *client) DeleteApplication(ctx context.Context, nameOrUUID string) error {
	resp, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", applicationsBasePathV1, nameOrUUID), nil, nil, nil)
	if err != nil {
		return err
	}

	return handleAPIErrors(resp)
}

// DeleteApplicationEntities does a bit more than DeleteApplication. Its main purpose is to remove underneath's app resources
// even in the event the application didn't exist.
func (c *client) DeleteApplicationEntities(ctx context.Context, nameOrUUID string) (*http.Response, error) {
	respAppDelete, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", applicationsBasePathV1, nameOrUUID), nil, nil, nil)
	if err != nil {
		return respAppDelete, err
	}

	// It is possible that an app failed to be created, but its resources still exist.
	if respAppDelete.StatusCode == 404 {
		respPipelineGet, err := c.GetPipelineByName(ctx, fmt.Sprintf("turbine-pipeline-%s", nameOrUUID))
		// If pipeline doesn't exist either, returns as if the app didn't exist in the first place
		if err != nil {
			return nil, handleAPIErrors(respAppDelete)
		}

		// Fetch connectors associated to that pipeline and delete each one.
		respConnectorsList, _ := c.ListPipelineConnectors(ctx, respPipelineGet.Name)

		// Delete destination connectors first
		destConnectors := filterConnectorsPerType(respConnectorsList, ConnectorTypeDestination)
		for _, connector := range destConnectors {
			_ = c.DeleteConnector(ctx, connector.Name)
		}

		// Delete source connectors
		srcConnectors := filterConnectorsPerType(respConnectorsList, ConnectorTypeSource)
		for _, connector := range srcConnectors {
			_ = c.DeleteConnector(ctx, connector.Name)
		}

		// Fetch all functions (we don't have way to filter functions from the API) and delete
		// the ones associated to the pipeline.
		respFunctionsList, _ := c.ListFunctions(ctx)
		for _, fn := range respFunctionsList {
			if fn.Pipeline.Name == respPipelineGet.Name {
				_, _ = c.DeleteFunction(ctx, fn.Name)
			}
		}

		// Delete pipeline as the last step
		err = c.DeletePipeline(ctx, respPipelineGet.Name)
		if err != nil {
			return nil, err
		}

		// Returns as if everything was successful
		resp := &http.Response{
			StatusCode: http.StatusNoContent,
		}
		return resp, handleAPIErrors(resp)
	}

	return respAppDelete, handleAPIErrors(respAppDelete)
}

func (c *client) GetApplication(ctx context.Context, nameOrUUID string) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", applicationsBasePathV1, nameOrUUID), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var a *Application
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c *client) ListApplications(ctx context.Context) ([]*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, applicationsBasePathV1, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var aa []*Application
	err = json.NewDecoder(resp.Body).Decode(&aa)
	if err != nil {
		return nil, err
	}

	return aa, nil
}

func (c *client) GetApplicationLogs(ctx context.Context, nameOrUUID string) (*ApplicationLogs, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/logs", applicationsBasePathV1, nameOrUUID), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var l *ApplicationLogs
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return nil, err
	}

	return l, nil
}
