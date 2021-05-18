package cased

import (
	"time"

	"github.com/dewski/jsonpath"
)

// Processors contains all processors available to transform an audit event
// before it's published to Cased.
var Processors = []Processor{
	SensitiveDataProcessor,
	PublishedAtProcessor,
}

// Processor is the interface necessary for processor functions to implement.
// It takes the audit event payload that is about to be published and mutates it
// as necessary.
//
// Each processor should be idempotent and not depend on another processor to be
// called beforehand.
type Processor func(*AuditEventPayload) *AuditEventPayload

// SensitiveDataProcessor adds sensitive data positions based on values.
func SensitiveDataProcessor(aep *AuditEventPayload) *AuditEventPayload {
	r := jsonpath.NewReader(aep.AuditEvent)
	for _, path := range r.Paths() {
		v := r.Path(path)
		if sv, ok := v.(SensitiveValue); ok {
			r := []*SensitiveRange{}
			for _, sr := range sv.Ranges {
				r = append(r, &sr)
			}

			aep.DotCased.PII[path] = r
		}
	}

	return aep
}

// PublishedAtProcessor sets the current time the audit event was published at.
func PublishedAtProcessor(aep *AuditEventPayload) *AuditEventPayload {
	aep.DotCased.PublishedAt = time.Now().UTC()

	return aep
}
