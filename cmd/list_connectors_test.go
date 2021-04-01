package cmd

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
)

func getConnectors() []*meroxa.Connector {
	var connectors []*meroxa.Connector
	c := utils.GenerateConnector()
	connectors = append(connectors, &c)
	return connectors
}

func TestListConnectorsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockListConnectorsClient(ctrl)

	connectors := getConnectors()

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

func TestListConnectorsOutput(t *testing.T) {
	flagRootOutputJSON = false
	connectors := getConnectors()

	output := utils.CaptureOutput(func() {
		ar := &ListConnectors{}
		ar.output(connectors)
	})

	connectorID := fmt.Sprintf("%d", connectors[0].ID)
	pipelineID := fmt.Sprintf("%d", connectors[0].PipelineID)

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
	flagRootOutputJSON = true
	connectors := getConnectors()

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
