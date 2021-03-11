package meroxa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"` // { "error" : "API error message" }
}

func handleAPIErrors(resp *http.Response) error {
	if resp.StatusCode > 204 {
		apiError, err := parseErrorFromBody(resp)

		// err if there was a problem decoding the resp.Body as the `errResponse` struct
		if err != nil {
			return err
		}

		// API error returned by Meroxa Platform API
		return apiError
	}
	return nil
}

func parseErrorFromBody(resp *http.Response) (error, error) {
	var er errResponse
	var body = resp.Body
	err := json.NewDecoder(body).Decode(&er)
	if err != nil {
		// In cases we didn't receive a proper JSON response
		if _, ok := err.(*json.SyntaxError); ok {
			return nil, errors.New(fmt.Sprintf("%s %s", resp.Proto, resp.Status))
		}

		return nil, err
	}

	return errors.New(er.Error), nil
}
