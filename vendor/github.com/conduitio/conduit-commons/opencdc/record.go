// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opencdc

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
)

// Record represents a single data record produced by a source and/or consumed
// by a destination connector.
// Record should be used as a value, not a pointer, except when (de)serializing
// the record. Note that methods related to (de)serializing the record mutate
// the record and are thus not thread-safe (see SetSerializer, FromProto and
// UnmarshalJSON).
type Record struct {
	// Position uniquely represents the record.
	Position Position `json:"position"`
	// Operation defines what triggered the creation of a record. There are four
	// possibilities: create, update, delete or snapshot. The first three
	// operations are encountered during normal CDC operation, while "snapshot"
	// is meant to represent records during an initial load. Depending on the
	// operation, the record will contain either the payload before the change,
	// after the change, both or none (see field Payload).
	Operation Operation `json:"operation"`
	// Metadata contains additional information regarding the record.
	Metadata Metadata `json:"metadata"`

	// Key represents a value that should identify the entity (e.g. database
	// row).
	Key Data `json:"key"`
	// Payload holds the payload change (data before and after the operation
	// occurred).
	Payload Change `json:"payload"`

	serializer RecordSerializer
}

// SetSerializer sets the serializer used to encode the record into bytes. If
// serializer is nil, the serializing behavior is reset to the default (JSON).
// This method mutates the receiver and is not thread-safe.
func (r *Record) SetSerializer(serializer RecordSerializer) {
	r.serializer = serializer
}

// Bytes returns the serialized representation of the Record. By default, this
// function returns a JSON representation. The serialization logic can be changed
// using SetSerializer.
func (r Record) Bytes() []byte {
	if r.serializer != nil {
		b, err := r.serializer.Serialize(r)
		if err == nil {
			return b
		}
		// serializer produced an error, fallback to default format
	}

	b, err := json.Marshal(r)
	if err != nil {
		// Unlikely to happen, records travel from/to plugins through GRPC.
		// If the content can be marshaled as protobuf it can be as JSON.
		panic(fmt.Errorf("error while marshaling Record as JSON: %w", err))
	}

	return b
}

func (r Record) Map() map[string]interface{} {
	var genericMetadata map[string]interface{}
	if r.Metadata != nil {
		genericMetadata = make(map[string]interface{}, len(r.Metadata))
		for k, v := range r.Metadata {
			genericMetadata[k] = v
		}
	}

	return map[string]any{
		"position":  []byte(r.Position),
		"operation": r.Operation.String(),
		"metadata":  genericMetadata,
		"key":       r.mapData(r.Key),
		"payload": map[string]interface{}{
			"before": r.mapData(r.Payload.Before),
			"after":  r.mapData(r.Payload.After),
		},
	}
}

func (r Record) mapData(d Data) interface{} {
	switch d := d.(type) {
	case StructuredData:
		return map[string]interface{}(d)
	case RawData:
		return []byte(d)
	}
	return nil
}

func (r Record) Clone() Record {
	var (
		metadata      map[string]string
		key           Data
		payloadBefore Data
		payloadAfter  Data
	)

	if r.Metadata != nil {
		metadata = make(map[string]string, len(r.Metadata))
		for k, v := range r.Metadata {
			metadata[k] = v
		}
	}

	if r.Key != nil {
		key = r.Key.Clone()
	}
	if r.Payload.Before != nil {
		payloadBefore = r.Payload.Before.Clone()
	}
	if r.Payload.After != nil {
		payloadAfter = r.Payload.After.Clone()
	}

	clone := Record{
		Position:  bytes.Clone(r.Position),
		Operation: r.Operation,
		Metadata:  metadata,
		Key:       key,
		Payload: Change{
			Before: payloadBefore,
			After:  payloadAfter,
		},
	}
	return clone
}
