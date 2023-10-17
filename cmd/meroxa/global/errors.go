package global

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type errResponse struct {
	Code    int                 `json:"code,omitempty"`
	Message string              `json:"message,omitempty"`
	Data    map[string][]string `json:"data,omitempty"`
}

func (err *errResponse) Error() string {
	msg := err.Message

	if errCount := len(err.Data); errCount > 0 {
		msg = fmt.Sprintf("%s. %d %s reported:%s",
			msg,
			errCount,
			func() string {
				if errCount > 1 {
					return "details"
				}
				return "detail"
			}(),
			mapToString(err.Data),
		)
	}
	return msg
}

func mapToString(m map[string][]string) string {
	s := ""
	count := 1

	// need to sort map keys separately
	var mKeys []string
	for k := range m {
		mKeys = append(mKeys, k)
	}
	sort.Strings(mKeys)

	for _, k := range mKeys {
		s = fmt.Sprintf("%s\n%d. %s: \"%s\"", s, count, k, strings.Join(m[k], `", "`))
		count++
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
	body := resp.Body
	err := json.NewDecoder(body).Decode(&er)
	if err != nil {
		// In cases we didn't receive a proper JSON response
		if _, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("%s %s", resp.Proto, resp.Status)
		}

		return nil, err
	}

	return &er, nil
}
