package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

const meroxaAccountUUIDHeader = "Meroxa-Account-UUID"
const accountsPath = "/v1/accounts"

type Account struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	CompanyName string `json:"companyName,omitempty"`
}

func (c *client) ListAccounts(ctx context.Context) ([]*Account, error) {
	path := accountsPath
	resp, err := c.MakeRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
