package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const functionsBasePath = "/v1/functions"

type Function struct {
	UUID         string             `json:"uuid"`
	Name         string             `json:"name"`
	InputStream  string             `json:"input_stream"`
	OutputStream string             `json:"output_stream"`
	Image        string             `json:"image"`
	Command      []string           `json:"command"`
	Args         []string           `json:"args"`
	EnvVars      map[string]string  `json:"env_vars"`
	Status       FunctionStatus     `json:"status"`
	Pipeline     PipelineIdentifier `json:"pipeline"`
	Logs         string             `json:"logs"` // CLI includes what's returned by GetFunctionLogs
}

type FunctionStatus struct {
	State   string `json:"state"`
	Details string `json:"details"`
}

type CreateFunctionInput struct {
	Name         string             `json:"name"`
	InputStream  string             `json:"input_stream"`
	OutputStream string             `json:"output_stream"`
	Pipeline     PipelineIdentifier `json:"pipeline"`
	Image        string             `json:"image"`
	Command      []string           `json:"command"`
	Args         []string           `json:"args"`
	EnvVars      map[string]string  `json:"env_vars"`
}

func (c *client) CreateFunction(ctx context.Context, input *CreateFunctionInput) (*Function, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, functionsBasePath, input, nil)
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

func (c *client) GetFunction(ctx context.Context, nameOrUUID string) (*Function, error) {
	path := fmt.Sprintf("%s/%s", functionsBasePath, nameOrUUID)

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

func (c *client) ListFunctions(ctx context.Context) ([]*Function, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, functionsBasePath, nil, nil)
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

func (c *client) DeleteFunction(ctx context.Context, nameOrUUID string) (*Function, error) {
	path := fmt.Sprintf("%s/%s", functionsBasePath, nameOrUUID)

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
