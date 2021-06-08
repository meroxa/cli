/*
Copyright © 2021 Meroxa Inc

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

package global

import (
	"fmt"
	"testing"
	"time"

	"github.com/meroxa/cli/utils"

	"github.com/spf13/viper"

	"github.com/cased/cased-go"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/spf13/cobra"
)

type metricsTestCase struct {
	cmd      *cobra.Command
	call     string
	rootArgs []string
	cmdArgs  []string
	want     map[string]interface{}
}

func cmdWithFlags(flags []string) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
	}

	for _, f := range flags {
		// we're assuming we're defining flags as bool just for testing purposes.
		cmd.Flags().Bool(f, true, "")
	}

	return cmd
}

func (tc *metricsTestCase) test(t *testing.T) {
	var called *cobra.Command
	run := func(c *cobra.Command, _ []string) {
		called = c
	}

	root := &cobra.Command{
		Use: "meroxa",
		Run: run,
	}

	tc.cmd.Run = run

	root.AddCommand(tc.cmd)
	root.SetArgs(tc.rootArgs)

	_ = root.Execute()

	if called == nil {
		if tc.call != "" {
			t.Errorf("missing expected call to command: %s", tc.call)
		}
		return
	}

	got := buildCommandInfo(called, tc.cmdArgs)

	if v := cmp.Diff(got, tc.want, cmpopts.IgnoreUnexported(cobra.Command{})); v != "" {
		t.Fatalf(v)
	}
}

func TestAddError(t *testing.T) {
	event := cased.AuditEvent{}

	addError(&event, nil)

	if v, ok := event["error"]; ok {
		t.Fatalf("not expected event to contain %q key, got %q", "error", v)
	}

	err := "unexpected error"
	addError(&event, fmt.Errorf(err))

	if v, ok := event["error"]; !ok || v != err {
		if !ok {
			t.Fatalf("expected event error to contain %q key", "error")
		}

		if v != err {
			t.Fatalf("expected event error to be %q, got %q", err, v)
		}
	}
}

func TestBuildCommandInfo(t *testing.T) {
	tests := map[string]metricsTestCase{
		"withAlias": {
			cmd: &cobra.Command{
				Use:     "list",
				Aliases: []string{"ls"},
			},
			call:     "list",
			rootArgs: []string{"ls"},
			want: map[string]interface{}{
				"alias": "ls",
			},
		},
		"withArgs": {
			cmd: &cobra.Command{
				Use: "list",
			},
			call:     "list",
			rootArgs: []string{"list"},
			cmdArgs:  []string{"arg1", "arg2"},
			want: map[string]interface{}{
				"args": "arg1,arg2",
			},
		},
		"withFlags": {
			cmd:      cmdWithFlags([]string{"flag1", "flag2"}),
			call:     "list",
			rootArgs: []string{"list", "--flag1", "--flag2"},
			want: map[string]interface{}{
				"flags": "flag1,flag2",
			},
		},
		"withDeprecated": {
			cmd: &cobra.Command{
				Use:        "list",
				Deprecated: "This command is deprecated",
			},
			call:     "list",
			rootArgs: []string{"list"},
			want: map[string]interface{}{
				"deprecated": true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, tc.test)
	}
}

func TestBuildBasicEvent(t *testing.T) {
	Version = "1.0.0"
	cmd := &cobra.Command{}

	want := cased.AuditEvent{
		"command":    map[string]interface{}{},
		"timestamp":  time.Now().UTC(),
		"user_agent": fmt.Sprintf("meroxa/%s darwin/amd64", Version),
	}

	got := buildBasicEvent(cmd, nil)
	if !cmp.Equal(want, got, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		if k == "timestamp" {
			if _, ok := v.(time.Time); !ok {
				t.Fatalf("expected %q to be %T got %T", k, time.Time{}, v)
			}
			return true
		}
		return false
	})) {
		t.Fatalf(cmp.Diff(got, want))
	}
}

func TestNewPublisherWithCasedAPIKey(t *testing.T) {
	Config = viper.New()
	defer clearConfiguration()

	apiKey := "8c32e3b7-d0e7-4650-a82b-e85e6a8d56fa"
	Config.Set("CASED_API_KEY", apiKey)

	got := NewPublisher()

	if got.Options().PublishKey != apiKey {
		t.Fatalf("expected publisher with API_KEY to be %q", apiKey)
	}
}

func TestNewPublisherWithoutCasedAPIKey(t *testing.T) {
	Config = viper.New()
	defer clearConfiguration()

	apiURL := fmt.Sprintf("%s/telemetry", meroxaBaseAPIURL)
	got := NewPublisher()

	if got.Options().PublishKey != "" {
		t.Fatalf("expected publisher without API_KEY set")
	}

	if got.Options().PublishURL != apiURL {
		t.Fatalf("expected publish url to be %q", apiURL)
	}
}

func TestNewPublisherPublishing(t *testing.T) {
	Config = viper.New()
	defer clearConfiguration()

	Config.Set("PUBLISH_METRICS", "false")
	got := NewPublisher()

	if !got.Options().Silence {
		t.Fatalf("expected publisher silence option to be %v", true)
	}

	Config.Set("PUBLISH_METRICS", "any-other-value")
	got = NewPublisher()

	if got.Options().Silence {
		t.Fatalf("expected publisher silence option to be %v", false)
	}
}

func TestNewPublisherWithDebug(t *testing.T) {
	Config = viper.New()
	defer clearConfiguration()

	Config.Set("CASED_DEBUG", true)
	got := NewPublisher()

	if !got.Options().Debug {
		t.Fatalf("expected publisher debug option to be %v", true)
	}
}

func clearConfiguration() {
	Config = nil
}

func TestPublishEventOnStdout(t *testing.T) {
	Config = viper.New()
	defer clearConfiguration()

	Config.Set("PUBLISH_METRICS", "stdout")

	event := cased.AuditEvent{
		"key": "event",
	}

	got := utils.CaptureOutput(func() {
		publishEvent(event)
	})

	want := "\n\nEvent: {\"key\":\"event\"}\n\n"

	if want != got {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestAddAction(t *testing.T) {
	root := &cobra.Command{Use: "meroxa"}

	list := &cobra.Command{Use: "list"}

	root.AddCommand(list)

	event := cased.AuditEvent{}
	addAction(&event, list)

	want := fmt.Sprintf("meroxa.%s", list.Use)

	if v := event["action"]; v != want {
		t.Fatalf("expected event action to be %q, got %q", want, v)
	}

	resources := &cobra.Command{Use: "resources"}

	list.AddCommand(resources)

	addAction(&event, resources)

	want = fmt.Sprintf("meroxa.%s.%s", list.Use, resources.Use)

	if v := event["action"]; v != want {
		t.Fatalf("expected event action to be %q, got %q", want, v)
	}
}
