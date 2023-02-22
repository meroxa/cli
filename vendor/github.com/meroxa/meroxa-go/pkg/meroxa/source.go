package meroxa

import (
	"context"
	"encoding/json"
	"net/http"
)

const sourcesBasePathV1 = "/v1/sources"
const sourcesBasePathV2 = "/v2/sources"

type Source struct {
	GetUrl string `json:"get_url"`
	PutUrl string `json:"put_url"`
}

type CreateSourceInputV2 struct {
	Environment *EntityIdentifier `json:"environment,omitempty"`
}

type SourceBlob struct {
	Url string `json:"url"`
}

// CreateSource using the v1 path won't accept body parameters, and it's an unauthenticated request
func (c *client) CreateSource(ctx context.Context) (*Source, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, sourcesBasePathV1, nil, nil, nil)
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

// CreateSourceV2 uses the v2 path, and this endpoint could receive an environment which will be used by Platform API
// to determine the data-plane where the source will be created.
func (c *client) CreateSourceV2(ctx context.Context, input *CreateSourceInputV2) (*Source, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, sourcesBasePathV2, input, nil, nil)
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
