package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/volatiletech/null/v8"
)

type ApplicationState string

const (
	ApplicationStateReady ApplicationState = "ready"
)

const applicationsBasePath = "/v1/applications"

// Application represents the Meroxa Application type within the Meroxa API
type Application struct {
	UUID      string               `json:"uuid"`
	Name      string               `json:"name"`
	Language  string               `json:"language"`
	Status    ApplicationStatus    `json:"status,omitempty"`
	Functions []FunctionIdentifier `json:"functions,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	DeletedAt time.Time            `json:"deleted_at,omitempty"`
}

type ApplicationIdentifier struct {
	UUID null.String `json:"uuid,omitempty"`
	Name null.String `json:"name,omitempty"`
}

// CreateApplicationInput represents the input for a Meroxa Application create operation in the API
type CreateApplicationInput struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type ApplicationStatus struct {
	State   ApplicationState `json:"state"`
	Details string           `json:"details,omitempty"`
}

func (c *client) CreateApplication(ctx context.Context, input *CreateApplicationInput) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, applicationsBasePath, input, nil)
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

func (c *client) DeleteApplication(ctx context.Context, name string) error {
	resp, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", applicationsBasePath, name), nil, nil)
	if err != nil {
		return err
	}

	return handleAPIErrors(resp)
}

func (c *client) GetApplication(ctx context.Context, name string) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", applicationsBasePath, name), nil, nil)
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
	resp, err := c.MakeRequest(ctx, http.MethodGet, applicationsBasePath, nil, nil)
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
