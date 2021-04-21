package meroxa

import (
	"context"
	"fmt"
	"net/http"
)

const connectorLogsBasePath = "/v1/connectors"

func (c *Client) GetConnectorLogs(ctx context.Context, connectorName string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/logs", connectorLogsBasePath, connectorName)

	req, err := c.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	// Override content-type and accept headers to text/palin
	req.Header.Add("Content-Type", textContentType)
	req.Header.Add("Accept", textContentType)

	return c.httpClient.Do(req)
}
