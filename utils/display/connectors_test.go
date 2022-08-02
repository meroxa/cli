package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/volatiletech/null/v8"
)

func TestConnectorRunningTable(t *testing.T) {
	connector := &meroxa.Connector{
		UUID:          "9483768f-c384-4b4a-96bf-b80a79a23b5c",
		Type:          "jdbc",
		Name:          "base",
		Configuration: nil,
		Metadata:      nil,
		Streams: map[string]interface{}{
			"dynamic": "false",
			"output":  []interface{}{"output-foo", "output-bar"},
		},
		State:        "running",
		Trace:        "",
		PipelineName: "pipeline-1",
		Environment:  &meroxa.EntityIdentifier{Name: null.StringFrom("my-env")},
	}
	failedConnector := &meroxa.Connector{}
	deepCopy(connector, failedConnector)
	failedConnector.State = "failed"
	failedConnector.Trace = "exception goes here"

	tests := map[string]*meroxa.Connector{
		"running": connector,
		"failed":  failedConnector,
	}

	tableHeaders := []string{"UUID", "ID", "Name", "Type", "Streams", "State", "Pipeline", "Environment"}

	for name, connector := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				fmt.Println(ConnectorTable(connector))
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "running":
				if !strings.Contains(out, connector.UUID) {
					t.Errorf("%s, not found", connector.UUID)
				}
				if !strings.Contains(out, connector.Name) {
					t.Errorf("%s, not found", connector.Name)
				}
				if !strings.Contains(out, connector.UUID) {
					t.Errorf("%s, not found", connector.UUID)
				}
			case "failed":
				if !strings.Contains(out, connector.UUID) {
					t.Errorf("%s, not found", connector.UUID)
				}
				if !strings.Contains(out, connector.Name) {
					t.Errorf("%s, not found", connector.Name)
				}
				if !strings.Contains(out, connector.UUID) {
					t.Errorf("%s, not found", connector.UUID)
				}
				if !strings.Contains(out, connector.Trace) {
					t.Errorf("%s, not found", connector.Trace)
				}
			}
			fmt.Println(out)
		})
	}
}

func TestConnectorsTable(t *testing.T) {
	connectionIDAlign := &meroxa.Connector{}
	connectionInputOutput := &meroxa.Connector{}
	connection := &meroxa.Connector{
		UUID:          "9483768f-c384-4b4a-96bf-b80a79a23b5c",
		Type:          "jdbc",
		Name:          "base",
		Configuration: nil,
		Metadata:      nil,
		Streams: map[string]interface{}{
			"dynamic": "false",
			"output":  []interface{}{"output-foo", "output-bar"},
		},
		State:        "running",
		Trace:        "",
		PipelineName: "pipeline-1",
		Environment:  &meroxa.EntityIdentifier{UUID: null.StringFrom("2c5326ac-041f-4679-b446-d6d95b91f497")},
	}

	deepCopy(connection, connectionIDAlign)
	connectionIDAlign.Name = "id-alignment"
	connectionIDAlign.UUID = "1000"

	deepCopy(connection, connectionInputOutput)
	connectionInputOutput.Name = "input-output"
	connectionInputOutput.Streams = map[string]interface{}{
		"dynamic": "true",
		"input":   []interface{}{"input-foo", "input-bar"},
		"output":  []interface{}{"output-foo", "output-bar"},
	}

	tests := map[string][]*meroxa.Connector{
		"Base":         {connection},
		"ID_Alignment": {connection, connectionIDAlign},
		"Input_Output": {connection, connectionInputOutput},
	}

	tableHeaders := []string{"UUID", "ID", "NAME", "TYPE", "STREAMS", "STATE", "PIPELINE", "ENVIRONMENT"}

	for name, connections := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintConnectorsTable(connections, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "Base":
				if !strings.Contains(out, connection.UUID) {
					t.Errorf("%s, not found", connection.UUID)
				}
				if !strings.Contains(out, connection.Name) {
					t.Errorf("%s, not found", connection.Name)
				}
				if !strings.Contains(out, connection.UUID) {
					t.Errorf("%s, not found", connection.UUID)
				}
			case "ID_Alignment":
				if !strings.Contains(out, connectionIDAlign.Name) {
					t.Errorf("%s, not found", connectionIDAlign.Name)
				}
				if !strings.Contains(out, connectionIDAlign.UUID) {
					t.Errorf("%s, not found", connectionIDAlign.UUID)
				}
			case "Input_Output":
				if !strings.Contains(out, connectionInputOutput.Name) {
					t.Errorf("%s, not found", connection.Name)
				}
				if !strings.Contains(out, "input-foo") {
					t.Errorf("%s, not found", "input-foo")
				}
				if !strings.Contains(out, "output-foo") {
					t.Errorf("%s, not found", "output-foo")
				}
			}
			fmt.Println(out)
		})
	}
}

func TestConnectorsTableWithoutHeaders(t *testing.T) {
	connection := &meroxa.Connector{
		UUID:          "9483768f-c384-4b4a-96bf-b80a79a23b5c",
		Type:          "jdbc",
		Name:          "base",
		Configuration: nil,
		Metadata:      nil,
		Streams: map[string]interface{}{
			"dynamic": "false",
			"output":  []interface{}{"output-foo", "output-bar"},
		},
		State:        "running",
		Trace:        "",
		PipelineName: "pipeline-1",
	}

	tableHeaders := []string{"UUID", "ID", "NAME", "TYPE", "STREAMS", "STATE", "PIPELINE", "ENVIRONMENT"}

	var connections []*meroxa.Connector
	connections = append(connections, connection)

	out := utils.CaptureOutput(func() {
		PrintConnectorsTable(connections, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}
	if !strings.Contains(out, connection.UUID) {
		t.Errorf("%s, not found", connection.UUID)
	}
	if !strings.Contains(out, connection.Name) {
		t.Errorf("%s, not found", connection.Name)
	}
	if !strings.Contains(out, connection.UUID) {
		t.Errorf("%s, not found", connection.UUID)
	}
}
