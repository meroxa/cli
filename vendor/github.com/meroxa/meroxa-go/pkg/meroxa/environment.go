package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const environmentsBasePath = "/v1/environments"

type EnvironmentState string

const (
	EnvironmentStateProvisioning   EnvironmentState = "provisioning"
	EnvironmentStateProvisioned    EnvironmentState = "provisioned"
	EnvironmentStateUpdating       EnvironmentState = "updating"
	EnvironmentStateError          EnvironmentState = "error"
	EnvironmentStateRepairing      EnvironmentState = "repairing"
	EnvironmentStateDeprovisioning EnvironmentState = "deprovisioning"
	EnvironmentStateDeprovisioned  EnvironmentState = "deprovisioned"
)

type EnvironmentViewStatus struct {
	State   EnvironmentState `json:"state"`
	Details string           `json:"details,omitempty"`
}

type EnvironmentRegion string

const (
	EnvironmentRegionAfSouth      EnvironmentRegion = "af-south-1"
	EnvironmentRegionApEast       EnvironmentRegion = "ap-east-1"
	EnvironmentRegionApNortheast1 EnvironmentRegion = "ap-northeast-1"
	EnvironmentRegionApNortheast2 EnvironmentRegion = "ap-northeast-2"
	EnvironmentRegionApNortheast3 EnvironmentRegion = "ap-northeast-3"
	EnvironmentRegionApSouth      EnvironmentRegion = "ap-south-1"
	EnvironmentRegionApSoutheast1 EnvironmentRegion = "ap-southeast-1"
	EnvironmentRegionApSoutheast2 EnvironmentRegion = "ap-southeast-2"
	EnvironmentRegionCaCentral    EnvironmentRegion = "ca-central-1"
	EnvironmentRegionEuCentral    EnvironmentRegion = "eu-central-1"
	EnvironmentRegionEuNorth      EnvironmentRegion = "eu-north-1"
	EnvironmentRegionEuSouth      EnvironmentRegion = "eu-south-1"
	EnvironmentRegionEuWest1      EnvironmentRegion = "eu-west-1"
	EnvironmentRegionEuWest2      EnvironmentRegion = "eu-west-2"
	EnvironmentRegionEuWest3      EnvironmentRegion = "eu-west-3"
	EnvironmentRegionMeSouth      EnvironmentRegion = "me-south-1"
	EnvironmentRegionSaEast1      EnvironmentRegion = "sa-east-1"
	EnvironmentRegionUsEast1      EnvironmentRegion = "us-east-1"
	EnvironmentRegionUsEast2      EnvironmentRegion = "us-east-2"
	EnvironmentRegionUsWest2      EnvironmentRegion = "us-west-2"
)

type EnvironmentType string

const (
	EnvironmentTypeSelfHosted EnvironmentType = "self_hosted"
	EnvironmentTypePrivate    EnvironmentType = "private"
	EnvironmentTypeCommon     EnvironmentType = "common"
)

type EnvironmentProvider string

const (
	EnvironmentProviderAws EnvironmentProvider = "aws"
)

// Environment represents the Meroxa Environment type within the Meroxa API
type Environment struct {
	UUID          string                 `json:"uuid"`
	Name          string                 `json:"name"`
	Provider      EnvironmentProvider    `json:"provider"`
	Region        EnvironmentRegion      `json:"region"`
	Type          EnvironmentType        `json:"type"`
	Configuration map[string]interface{} `json:"config,omitempty"`
	Status        EnvironmentViewStatus  `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// CreateEnvironmentInput represents the input for a Meroxa Environment we're creating within the Meroxa API
type CreateEnvironmentInput struct {
	Type          EnvironmentType        `json:"type,omitempty"`
	Provider      EnvironmentProvider    `json:"provider,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Configuration map[string]interface{} `json:"config"`
	Region        EnvironmentRegion      `json:"region,omitempty"`
}

// ListEnvironments returns an array of Environments (scoped to the calling user)
func (c *client) ListEnvironments(ctx context.Context) ([]*Environment, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, environmentsBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var ee []*Environment
	err = json.NewDecoder(resp.Body).Decode(&ee)
	if err != nil {
		return nil, err
	}

	return ee, nil
}

// CreateEnvironment creates a new Environment based on a CreateEnvironmentInput
func (c *client) CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*Environment, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, environmentsBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var e Environment
	err = json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (c *client) GetEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error) {
	path := fmt.Sprintf("%s/%s", environmentsBasePath, nameOrUUID)
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var e *Environment
	err = json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (c *client) DeleteEnvironment(ctx context.Context, nameOrUUID string) (*Environment, error) {
	path := fmt.Sprintf("%s/%s", environmentsBasePath, nameOrUUID)
	resp, err := c.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var e *Environment
	err = json.NewDecoder(resp.Body).Decode(&e)

	return e, nil
}
