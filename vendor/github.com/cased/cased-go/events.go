package cased

import (
	"encoding/json"
	"time"
)

// EventPayload is the JSON event.
type EventPayload map[string]interface{}

type Event struct {
	// The Event ID
	ID string `json:"id"`

	// The API URL for the workflow.
	APIURL string `json:"api_url"`

	// Result contains information about the workflow run. If a workflow was not
	// specified or detected, it will be empty.
	Result Result `json:"result"`

	// Event has been processed.
	Event EventPayload `json:"event"`

	// OriginalEvent contains the original event published to Cased.
	OriginalEvent EventPayload `json:"original_event"`

	// UpdatedAt is when the event was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the event was published.
	CreatedAt time.Time `json:"created_at"`
}

// EventParams contains the available fields when publishing events.
type EventParams struct {
	Params

	// WorkflowID is optional and only required if the workflow is known ahead of
	// time.
	WorkflowID *string

	// Event is the event that is published to Cased.
	Event EventPayload
}

func (ep EventParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(ep.Event)
}
