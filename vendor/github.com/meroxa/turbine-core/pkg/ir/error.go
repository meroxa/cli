package ir

import (
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type SpecValidationError struct {
	*jsonschema.ValidationError
	message string
}

func (e *SpecValidationError) Error() string {
	return e.message
}

func rootError(err *jsonschema.ValidationError) *jsonschema.ValidationError {
	if len(err.Causes) > 0 {
		return rootError(err.Causes[0])
	}
	return err
}

func NewSpecValidationError(err error) error {
	var e *jsonschema.ValidationError
	if errors.As(err, &e) {
		rootErr := rootError(e)
		return &SpecValidationError{
			ValidationError: e,
			message:         fmt.Sprintf("%q field fails %s validation: %s", rootErr.InstanceLocation, rootErr.KeywordLocation, rootErr.Message),
		}
	}

	return &SpecValidationError{
		message: err.Error(),
	}
}
