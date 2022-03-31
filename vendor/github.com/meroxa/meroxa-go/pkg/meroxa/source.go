package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

const sourcesBasePath = "/v1/sources"

type Source struct {
	GetUrl string `json:"get_url"`
	PutUrl string `json:"put_url"`
}

type SourceBlob struct {
	Url string `json:"url"`
}

func (c *client) CreateSource(ctx context.Context) (*Source, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, sourcesBasePath, nil, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var s *Source
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
