package cased

import "encoding/json"

// SensitiveValue contains the sensitive value and all the sensitive ranges
// within the provided value.
type SensitiveValue struct {
	Value  string
	Ranges []SensitiveRange
}

// MarshalJSON encodes the provided sensitive value for JSON representation.
func (sv SensitiveValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(sv.Value)
}

// SensitiveRange is a range that informs Cased about any sensitive information
// stored in an AuditEvent.
type SensitiveRange struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	Label string `json:"label"`
}

// NewSensitiveValue marks an entire string as sensitive.
//
// The marked sensitive value will be encoded upon publishing to Cased.
func NewSensitiveValue(value, label string) SensitiveValue {
	return SensitiveValue{
		Value: value,
		Ranges: []SensitiveRange{
			{
				Begin: 0,
				End:   len(value),
				Label: label,
			},
		},
	}
}
