package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/volatiletech/null/v8"
)

type Function struct {
	UUID         string                `json:"uuid"`
	Name         string                `json:"name"`
	InputStream  string                `json:"input_stream"`
	OutputStream string                `json:"output_stream"`
	Image        string                `json:"image"`
	Command      []string              `json:"command"`
	Args         []string              `json:"args"`
	EnvVars      map[string]string     `json:"env_vars"`
	Status       FunctionStatus        `json:"status"`
	Pipeline     PipelineIdentifier    `json:"pipeline"`
	Application  ApplicationIdentifier `json:"application,omitempty"`
}

type FunctionStatus struct {
	State   string `json:"state"`
	Details string `json:"details"`
}

type FunctionIdentifier struct {
	Name null.String `json:"name,omitempty"`
	UUID null.String `json:"uuid,omitempty"`
}

type CreateFunctionInput struct {
	Name         string                `json:"name"`
	InputStream  string                `json:"input_stream"`
	OutputStream string                `json:"output_stream"`
	Pipeline     PipelineIdentifier    `json:"pipeline"`
	Application  ApplicationIdentifier `json:"application"`
	Image        string                `json:"image"`
	Command      []string              `json:"command"`
	Args         []string              `json:"args"`
	EnvVars      map[string]string     `json:"env_vars"`
}

func functionsPath(appNameOrUUID, nameOrUUID string) string {
	path := fmt.Sprintf("%s/%s/functions", applicationsBasePath, appNameOrUUID)
	if nameOrUUID != "" {
		path += fmt.Sprintf("/%s", nameOrUUID)
	}
	return path
}

func (c *client) CreateFunction(ctx context.Context, input *CreateFunctionInput) (*Function, error) {
	var appID string
	if input.Application.Name.Valid && input.Application.Name.String != "" {
		appID = input.Application.Name.String
	} else if input.Application.UUID.Valid && input.Application.UUID.String != "" {
		appID = input.Application.UUID.String
	}
	if appID == "" {
		return nil, fmt.Errorf("application identifier not provided")
	}
	path := functionsPath(appID, "")
	resp, err := c.MakeRequest(ctx, http.MethodPost, path, input, nil)
	if err != nil {
		return nil, err
	}

	if err := handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var fun Function
	if err := json.NewDecoder(resp.Body).Decode(&fun); err != nil {
		return nil, err
	}

	return &fun, nil
}

func (c *client) GetFunction(ctx context.Context, appNameOrUUID, nameOrUUID string) (*Function, error) {
	path := functionsPath(appNameOrUUID, nameOrUUID)

	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var fun Function
	err = json.NewDecoder(resp.Body).Decode(&fun)
	if err != nil {
		return nil, err
	}

	return &fun, nil
}

func (c *client) ListFunctions(ctx context.Context, appNameOrUUID string) ([]*Function, error) {
	path := functionsPath(appNameOrUUID, "")
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var funs []*Function
	err = json.NewDecoder(resp.Body).Decode(&funs)
	if err != nil {
		return nil, err
	}

	return funs, nil
}

func (c *client) DeleteFunction(ctx context.Context, appNameOrUUID, nameOrUUID string) (*Function, error) {
	path := functionsPath(appNameOrUUID, nameOrUUID)

	resp, err := c.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var fun Function
	err = json.NewDecoder(resp.Body).Decode(&fun)
	if err != nil {
		return nil, err
	}

	return &fun, nil
}
