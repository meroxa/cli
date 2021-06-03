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
	"testing"

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
		t.Logf("called: %q", c.Name())
		t.Logf("called as: %q", c.CalledAs())
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
