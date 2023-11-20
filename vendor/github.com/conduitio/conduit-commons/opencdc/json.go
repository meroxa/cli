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
	"fmt"

	"github.com/goccy/go-json"
)

func (r *Record) UnmarshalJSON(b []byte) error {
	var raw struct {
		Position  Position  `json:"position"`
		Operation Operation `json:"operation"`
		Metadata  Metadata  `json:"metadata"`
		Payload   struct {
			Before json.RawMessage `json:"before"`
			After  json.RawMessage `json:"after"`
		} `json:"payload"`
		Key json.RawMessage `json:"key"`
	}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err //nolint:wrapcheck // no additional context to add
	}

	key, err := dataUnmarshalJSON(raw.Key)
	if err != nil {
		return err
	}

	payloadBefore, err := dataUnmarshalJSON(raw.Payload.Before)
	if err != nil {
		return err
	}

	payloadAfter, err := dataUnmarshalJSON(raw.Payload.After)
	if err != nil {
		return err
	}

	r.Position = raw.Position
	r.Operation = raw.Operation
	r.Metadata = raw.Metadata
	r.Key = key
	r.Payload = Change{
		Before: payloadBefore,
		After:  payloadAfter,
	}

	return nil
}

func dataUnmarshalJSON(b []byte) (Data, error) {
	if b[0] == '"' {
		var data RawData
		err := json.Unmarshal(b, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal raw data: %w", err)
		}
		return data, nil
	}
	var data StructuredData
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal structured data: %w", err)
	}
	return data, nil
}
