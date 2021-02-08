package meroxa

import (
	"context"
	"fmt"
	"net/http"
)

const connectorLogsBasePath = "/v1/connectors"

func (c *Client) GetConnectorLogs(ctx context.Context, connectorName string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", connectorLogsBasePath, connectorName)
	return c.makeRequest(ctx, http.MethodGet, path, nil, nil)
}
