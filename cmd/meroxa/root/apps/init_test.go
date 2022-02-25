package apps

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestInitAppArgs(t *testing.T) {
	tests := []struct {
		args    []string
		err     error
		appName string
	}{
		{args: nil, err: errors.New("requires an application name"), appName: ""},
		{args: []string{"my-app-name"}, err: nil, appName: "my-app-name"},
	}

	for _, tt := range tests {
		cc := &Init{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.appName != cc.args.appName {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.appName, cc.args.appName)
		}
	}
}

func TestInitAppFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "lang", shorthand: "l", required: true},
		{name: "path", required: false},
	}

	c := builder.BuildCobraCommand(&Init{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}

		if cf.Hidden != f.hidden {
			if cf.Hidden {
				t.Fatalf("expected flag \"%s\" not to be hidden", f.name)
			} else {
				t.Fatalf("expected flag \"%s\" to be hidden", f.name)
			}
		}
	}
}

func TestGitInit(t *testing.T) {
	testDir := os.TempDir() + "/tests" + uuid.New().String()

	tests := []struct {
		path string
		err  error
	}{
		{path: "", err: errors.New("path is required")},
		{path: testDir, err: nil},
	}

	for _, tt := range tests {
		cc := &Init{}
		err := cc.GitInit(context.Background(), tt.path)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.err == nil {
			if _, err := os.Stat(testDir + "/.git"); os.IsNotExist(err) {
				t.Fatalf("expected directory \"%s\" to be created", testDir)
			}
		}
	}

	os.RemoveAll(testDir)
}
