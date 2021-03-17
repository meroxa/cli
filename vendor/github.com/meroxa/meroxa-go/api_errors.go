package meroxa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type errResponse struct {
	ErrorDeprecated string              `json:"error,omitempty"` // { "error" : "API error message" }
	Code            string              `json:"code,omitempty"`
	Message         string              `json:"message,omitempty"`
	Details         map[string][]string `json:"details,omitempty"`
}

func (err *errResponse) Error() string {
	if err.ErrorDeprecated != "" {
		return err.ErrorDeprecated
	}

	msg := fmt.Sprintf("\ncode: %q\nmessage: %q",
		err.Code,
		err.Message,
	)
	if len(err.Details) > 0 {
		msg = fmt.Sprintf("%s\ndetails: %s",
			msg,
			mapToString(err.Details),
		)
	}
	return msg
}

func mapToString(m map[string][]string) string {
	s := ""
	for k, v := range m {
		s = fmt.Sprintf("%s\n\t%q: [\"%s\"]", s, k, strings.Join(v, `", "`))
	}
	return s
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

	return &er, nil
}
