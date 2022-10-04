package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeploymentState string

const (
	DeploymentStateDeploying        DeploymentState = "deploying"
	DeploymentStateDeployingError   DeploymentState = "deploying_error"
	DeploymentStateRollingBack      DeploymentState = "rolling_back"
	DeploymentStateRollingBackError DeploymentState = "rolling_back_error"
	DeploymentStateDeployed         DeploymentState = "deployed"
)

type DeploymentStatus struct {
	State   DeploymentState `json:"state"`
	Details string          `json:"details,omitempty"`
}

type Deployment struct {
	UUID        string                 `json:"uuid"`
	GitSha      string                 `json:"git_sha"`
	Application EntityIdentifier       `json:"application"`
	CreatedAt   time.Time              `json:"created_at"`
	DeletedAt   time.Time              `json:"deleted_at,omitempty"`
	Status      DeploymentStatus       `json:"status"`
	Spec        map[string]interface{} `json:"spec,omitempty"`
	SpecVersion string                 `json:"spec_version,omitempty"`
	CreatedBy   string                 `json:"created_by"`
}

type CreateDeploymentInput struct {
	GitSha      string                 `json:"git_sha"`
	Application EntityIdentifier       `json:"application"`
	Spec        map[string]interface{} `json:"spec,omitempty"`
	SpecVersion string                 `json:"spec_version,omitempty"`
}

func (c *client) GetLatestDeployment(ctx context.Context, appIdentifier string) (*Deployment, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/deployments/latest", applicationsBasePath, appIdentifier), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var d *Deployment
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}

func (c *client) GetDeployment(ctx context.Context, appIdentifier string, depUUID string) (*Deployment, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s/deployments/%s", applicationsBasePath, appIdentifier, depUUID), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var d *Deployment
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}

func (c *client) CreateDeployment(ctx context.Context, input *CreateDeploymentInput) (*Deployment, error) {
	appIdentifier, err := input.Application.GetNameOrUUID()

	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(ctx, http.MethodPost, fmt.Sprintf("%s/%s/deployments", applicationsBasePath, appIdentifier), input, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var d *Deployment
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}
