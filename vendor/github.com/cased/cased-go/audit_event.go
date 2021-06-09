package cased

import (
	"encoding/json"
	"time"
)

const (
	// DotCasedKey is the key name for the property that encodes information about
	// an AuditEvent.
	DotCasedKey = ".cased"

	// DefaultSensitiveLabel is the default value used if a particular PII
	// does not contain a label.
	DefaultSensitiveLabel = "sensitive-value"
)

// DotCased is a reserved property in an audit event containing the original
// event, any modifications to the event post-processing, timestamps, and more.
type DotCased struct {
	PII                map[string][]*SensitiveRange `json:"pii,omitempty"`
	ID                 string                       `json:"id,omitempty"`
	Event              AuditEvent                   `json:"event,omitempty"`
	PublisherUserAgent string                       `json:"publisher_user_agent,omitempty"`
	ProcessedAt        *time.Time                   `json:"processed_at,omitempty"`
	ReceivedAt         *time.Time                   `json:"received_at,omitempty"`
	PublishedAt        time.Time                    `json:"published_at"`
}

// AuditEvent ...
type AuditEvent map[string]interface{}

// MarshalJSON ...
func (ae AuditEvent) MarshalJSON() ([]byte, error) {
	// Need to alias AuditEvent so we do not have an infinite loop when calling
	// MarshalJSON.
	type AE AuditEvent

	return json.Marshal(AE(ae))
}

// NewAuditEventPayload ...
func NewAuditEventPayload(event AuditEvent) *AuditEventPayload {
	aep := &AuditEventPayload{
		DotCased: DotCased{
			PII: map[string][]*SensitiveRange{},
		},
		AuditEvent: event,
	}

	aep.process()

	return aep
}

// AuditEventPayload is the wrapper struct hosting the nestable JSON AuditEvent
// with the internal `.cased` property with a rich struct.
type AuditEventPayload struct {
	DotCased   DotCased `json:".cased"`
	AuditEvent AuditEvent
}

// MarshalJSON ...
func (aep *AuditEventPayload) MarshalJSON() ([]byte, error) {
	f := map[string]interface{}{}
	for k, v := range aep.AuditEvent {
		f[k] = v
	}

	f[DotCasedKey] = aep.DotCased
	return json.Marshal(f)
}

// UnmarshalJSON ...
func (aep *AuditEventPayload) UnmarshalJSON(data []byte) error {
	// We cannot use AuditEventPayload as it doesn't know how to serialize the
	// other properties into AuditEvent. This inline struct enables us to only
	// extract the `.cased` property.
	var dc struct {
		DotCased DotCased `json:".cased"`
	}
	if err := json.Unmarshal(data, &dc); err != nil {
		return err
	}
	aep.DotCased = dc.DotCased

	// Serialize all attributes to a common interface.
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ae := AuditEvent{}
	for key, value := range v {
		// We've already extracted out the `.cased` property in the beginning of
		// UnmarshalJSON so we need to ignore it here.
		if key == DotCasedKey {
			continue
		}

		ae[key] = value
	}
	aep.AuditEvent = ae

	return nil
}

func (aep *AuditEventPayload) process() {
	for _, processor := range Processors {
		processor(aep)
	}
}
