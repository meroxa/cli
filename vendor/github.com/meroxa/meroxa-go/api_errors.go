package meroxa

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func handleAPIErrors(resp *http.Response) error {
	if resp.StatusCode > 204 {
		apiError, err := parseErrorFromBody(resp.Body)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func parseErrorFromBody(body io.ReadCloser) (error, error) {
	var er errResponse
	err := json.NewDecoder(body).Decode(&er)
	if err != nil {
		return nil, err
	}

	return errors.New(er.Error), nil
}
