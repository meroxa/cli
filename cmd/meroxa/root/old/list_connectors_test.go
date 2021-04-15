package old

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func getConnectors(pipelineID int) []*meroxa.Connector {
	var connectors []*meroxa.Connector
	c := utils.GenerateConnector(pipelineID)
	connectors = append(connectors, &c)
	return connectors
}

func TestListConnectorsFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"pipeline", false, ""},
	}

	c := &cobra.Command{}
	lc := &ListConnectors{}
	lc.setFlags(c)

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

func TestListConnectorsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListConnectorsClient(ctrl)

	connectors := getConnectors(1)

	client.
		EXPECT().
		ListConnectors(ctx).
		Return(connectors, nil)

	lc := &ListConnectors{}
	got, err := lc.execute(ctx, client)

	if !reflect.DeepEqual(got, connectors) {
		t.Fatalf("expected \"%v\", got \"%v\"", connectors, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestListPipelineConnectorsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListConnectorsClient(ctrl)

	pipeline := utils.GeneratePipeline()
	connectors := getConnectors(pipeline.ID)

	lc := &ListConnectors{}
	lc.pipeline = pipeline.Name

	client.
		EXPECT().
		GetPipelineByName(ctx, lc.pipeline).
		Return(&pipeline, nil)

	client.
		EXPECT().
		ListPipelineConnectors(ctx, pipeline.ID).
		Return(connectors, nil)

	got, err := lc.execute(ctx, client)

	if !reflect.DeepEqual(got, connectors) {
		t.Fatalf("expected \"%v\", got \"%v\"", connectors, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestListConnectorsOutput(t *testing.T) {
	FlagRootOutputJSON = false
	connectors := getConnectors(1)

	output := utils.CaptureOutput(func() {
		ar := &ListConnectors{}
		ar.output(connectors)
	})

	pipelineID := fmt.Sprintf("%d", connectors[0].PipelineID)
	connectorID := fmt.Sprintf("%d", connectors[0].ID)

	if !strings.Contains(output, connectorID) {
		t.Fatalf("expected output \"%s\" got \"%s\"", connectorID, output)
	}

	if !strings.Contains(output, connectors[0].Name) {
		t.Fatalf("expected output \"%s\" got \"%s\"", connectors[0].Name, output)
	}

	if !strings.Contains(output, connectors[0].Type) {
		t.Fatalf("expected output \"%s\" got \"%s\"", connectors[0].Type, output)
	}

	if !strings.Contains(output, pipelineID) {
		t.Fatalf("expected output \"%s\" got \"%s\"", pipelineID, output)
	}
}

func TestListConnectorsJSONOutput(t *testing.T) {
	FlagRootOutputJSON = true
	connectors := getConnectors(1)

	output := utils.CaptureOutput(func() {
		ar := &ListConnectors{}
		ar.output(connectors)
	})

	var parsedOutput []*meroxa.Connector
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(connectors, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
