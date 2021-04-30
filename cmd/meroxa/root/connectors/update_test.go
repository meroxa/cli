package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/meroxa/meroxa-go"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/utils"
)

func TestUpdateConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires connector name"), name: ""},
		{args: []string{"conName"}, err: nil, name: "conName"},
	}

	for _, tt := range tests {
		cc := &UpdateConnector{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestUpdateConnectorFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
		hidden    bool
	}{
		{name: "state", required: true, shorthand: "", hidden: false},
	}

	c := builder.BuildCobraCommand(&UpdateConnector{})

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

func TestUpdateConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockUpdateConnectorClient(ctrl)
	logger := log.NewTestLogger()

	u := &UpdateConnector{
		client: client,
		logger: logger,
	}

	c := utils.GenerateConnector(0, "")
	u.args.Name = c.Name
	u.flags.State = "pause"

	client.
		EXPECT().
		UpdateConnectorStatus(ctx, u.args.Name, u.flags.State).
		Return(&c, nil)

	err := u.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Updating connector %q...
Connector %q successfully updated!
`, u.args.Name, u.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, c) {
		t.Fatalf("expected \"%v\", got \"%v\"", c, gotConnector)
	}
}
