package meroxa

import (
	"context"
	"fmt"
	"net/http"
)

const (
	connectorLogsBasePath = "/v1/connectors"
)

func (c *client) GetConnectorLogs(ctx context.Context, nameOrID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", connectorLogsBasePath, nameOrID)

	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	// Override content-type and accept headers to text/palin
	req.Header.Add("Content-Type", textContentType)
	req.Header.Add("Accept", textContentType)

	return c.httpClient.Do(req)
}

func (c *client) GetFunctionLogs(ctx context.Context, appNameOrUUID, nameOrUUID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/functions/%s/logs", applicationsBasePath, appNameOrUUID, nameOrUUID)

	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	// Override content-type and accept headers to text/palin
	req.Header.Add("Content-Type", textContentType)
	req.Header.Add("Accept", textContentType)

	return c.httpClient.Do(req)
}
