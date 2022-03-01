/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliee.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/spf13/viper"

	"github.com/meroxa/cli/log"
)

func TestConfigDescribeExecution(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	type confEnv struct {
		Path   string                 `json:"path"`
		Config map[string]interface{} `json:"config"`
	}

	var env confEnv

	env.Path = "~/meroxa"
	env.Config = make(map[string]interface{})
	env.Config["my-key"] = "my-value"
	env.Config["my-token"] = "supersecrettokenandverylongaswell"

	cfg := viper.New()
	cfg.Set("my-key", env.Config["my-key"])
	cfg.Set("my-token", env.Config["my-token"])

	cfg.SetConfigFile(env.Path)
	global.Config = cfg

	e := &Describe{
		logger: logger,
	}

	err := e.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Using meroxa config located in "%s"

my-key: %s
my-token: superse...gaswell
`, env.Path, env.Config["my-key"])

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()

	var gotEnv confEnv
	err = json.Unmarshal([]byte(gotJSONOutput), &gotEnv)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(env, gotEnv) {
		t.Fatalf("expected \"%v\", got \"%v\"", gotEnv, env)
	}
}
