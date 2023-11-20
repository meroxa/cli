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
	"errors"
)

var (
	// ErrMetadataFieldNotFound is returned in metadata utility functions when a
	// metadata field is not found.
	ErrMetadataFieldNotFound = errors.New("metadata field not found")
	// ErrUnknownOperation is returned when trying to parse an Operation string
	// and encountering an unknown operation.
	ErrUnknownOperation = errors.New("unknown operation")

	// ErrInvalidProtoDataType is returned when trying to convert a proto data
	// type to raw or structured data.
	ErrInvalidProtoDataType = errors.New("invalid proto data type")
)
