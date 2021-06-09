package cased

import "encoding/json"

type ErrorCode string

const (
	ErrorCodeNotFound                    ErrorCode = "not_found"
	ErrorCodeInvalidContentType          ErrorCode = "invalid_content_type"
	ErrorCodeInvalidAuthenticationScheme ErrorCode = "invalid_authentication_scheme"
	ErrorCodeReadOnlyAPIKey              ErrorCode = "read_only_api_key"
	ErrorCodeInvalidAPIKey               ErrorCode = "invalid_api_key"
	ErrorCodeUnauthorized                ErrorCode = "unauthorized"
	ErrorCodeInvalidRequest              ErrorCode = "invalid_request"
)

type ErrorMessage struct {
	Resource string `json:"resource"`
	Path     string `json:"path"`
	Code     string `json:"code"`
}

type Error struct {
	Code    ErrorCode       `json:"error,omitempty"`
	Message string          `json:"message"`
	Errors  []*ErrorMessage `json:"errors,omitempty"`

	Err error `json:"-"`
}

func (ae *Error) Error() string {
	data, _ := json.Marshal(ae)
	return string(data)
}
