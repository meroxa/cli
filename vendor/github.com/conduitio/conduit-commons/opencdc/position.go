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

// Position is a unique identifier for a record being processed.
// It's a Source's responsibility to choose and assign record positions,
// as they will be used by the Source in subsequent pipeline runs.
type Position []byte

// String is used when displaying the position in logs.
func (p Position) String() string {
	if p != nil {
		return string(p)
	}
	return "<nil>"
}
