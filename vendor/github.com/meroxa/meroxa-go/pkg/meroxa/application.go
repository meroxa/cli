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
	ApplicationStateReady ApplicationState = "ready"
)

const applicationsBasePath = "/v1/applications"

// Application represents the Meroxa Application type within the Meroxa API
type Application struct {
	UUID      string            `json:"uuid"`
	Name      string            `json:"name"`
	Language  string            `json:"language"`
	Status    ApplicationStatus `json:"status,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// CreateApplicationInput represents the input for a Meroxa Application create operation in the API
type CreateApplicationInput struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

// UpdateApplicationInput represents the input for a Meroxa Application update operation in the API
type UpdateApplicationInput struct {
	Name string `json:"name"`
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

func (c *client) UpdateApplication(ctx context.Context, uuid string, input *UpdateApplicationInput) (*Application, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", applicationsBasePath, uuid), input, nil)
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
