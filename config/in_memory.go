/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

type InMemoryConfig struct {
	values map[string]interface{}
}

// NewInMemoryConfig is the constructor for InMemoryConfig.
func NewInMemoryConfig() *InMemoryConfig {
	return &InMemoryConfig{
		values: make(map[string]interface{}),
	}
}

func (m *InMemoryConfig) Set(key string, value interface{}) {
	m.values[key] = value
}

func (m *InMemoryConfig) GetString(key string) string {
	return m.values[key].(string)
}
