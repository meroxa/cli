package deprecated

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
)

func TestCreateConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"conName"}, nil, "conName"},
	}

	for _, tt := range tests {
		cc := &CreateConnector{source: "source or destination is required"}
		err := cc.setArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.name)
		}
	}
}

func TestCreateConnectorFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"input", true, ""},
		{"config", false, "c"},
		{"from", false, ""},
		{"to", false, ""},
		{"pipeline", false, ""},
		{"metadata", false, "m"},
	}

	c := &cobra.Command{}
	cc := &CreateConnector{}
	cc.setFlags(c)

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
	}
}

func TestCreateConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockCreateConnectorClient(ctrl)

	cc := &CreateConnector{
		input:        "foo",
		config:       `{"key":"value"}`,
		metadata:     `{"metakey":"metavalue"}`,
		source:       "my-resource",
		destination:  "",
		name:         "connector-name",
		pipelineName: "my-pipeline",
	}

	cr := utils.GenerateConnector(0)

	client.
		EXPECT().
		GetResourceByName(
			ctx,
			"my-resource",
		).
		Return(&meroxa.Resource{ID: 123}, nil)

	client.
		EXPECT().
		CreateConnector(
			ctx,
			meroxa.CreateConnectorInput{
				Name:         "connector-name",
				ResourceID:   123,
				PipelineName: "my-pipeline",
				Configuration: map[string]interface{}{
					"key":   "value",
					"input": "foo",
				},
				Metadata: map[string]interface{}{
					"metakey":          "metavalue",
					"mx:connectorType": "source",
				},
			},
		).
		Return(&cr, nil)

	output := utils.CaptureOutput(func() {
		got, err := cc.execute(ctx, client)

		if !reflect.DeepEqual(got, &cr) {
			t.Fatalf("expected \"%v\", got \"%v\"", &cr, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := "Creating connector from source my-resource..."
	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestCreateConnectorOutput(t *testing.T) {
	c := utils.GenerateConnector(0)
	FlagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		cc := &CreateConnector{}
		cc.output(&c)
	})

	expected := fmt.Sprintf("Connector %s successfully created!", c.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestCreateConnectorJSONOutput(t *testing.T) {
	c := utils.GenerateConnector(0)
	FlagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		cc := &CreateConnector{}
		cc.output(&c)
	})

	var parsedOutput meroxa.Connector
	_ = json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(c, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
