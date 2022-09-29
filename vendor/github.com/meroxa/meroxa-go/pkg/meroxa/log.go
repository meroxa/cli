package meroxa

import (
	"context"
	"fmt"
	"net/http"
)

const (
	connectorLogsBasePath = "/v1/connectors"
	functionLogsBasePath  = "/v1/functions"
	buildLogsBasePath     = "/v1/builds"
)

func (c *client) GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", connectorLogsBasePath, nameOrID)
	headers := make(http.Header)
	headers["Content-Type"] = []string{textContentType}
	headers["Accept"] = []string{textContentType}
	return c.MakeRequest(ctx, http.MethodGet, path, nil, nil, headers)
}

func (c *client) GetFunctionLogs(ctx context.Context, nameOrUUID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", functionLogsBasePath, nameOrUUID)
	headers := make(http.Header)
	headers["Content-Type"] = []string{textContentType}
	headers["Accept"] = []string{textContentType}
	return c.MakeRequest(ctx, http.MethodGet, path, nil, nil, headers)
}

func (c *client) GetBuildLogs(ctx context.Context, uuid string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", buildLogsBasePath, uuid)
	headers := make(http.Header)
	headers["Content-Type"] = []string{textContentType}
	headers["Accept"] = []string{textContentType}
	return c.MakeRequest(ctx, http.MethodGet, path, nil, nil, headers)
}
