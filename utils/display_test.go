package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/meroxa/meroxa-go"
)

func TestResourcesTable(t *testing.T) {
	resource := &meroxa.Resource{
		ID:          0,
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
	}
	resIDAlign := &meroxa.Resource{
		ID:          1000,
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
	}

	tests := map[string][]*meroxa.Resource{
		"Base":         {resource},
		"ID_Alignment": {resource, resIDAlign},
	}

	for name, resources := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintResourcesTable(resources)
			})

			switch name {
			case "Base":
				if !strings.Contains(out, resource.Name) {
					t.Errorf("%s, not found", resource.Name)
				}
				if !strings.Contains(out, strconv.Itoa(resource.ID)) {
					t.Errorf("%d, not found", resource.ID)
				}
			case "ID_Alignment":
				if !strings.Contains(out, resIDAlign.Name) {
					t.Errorf("%s, not found", resIDAlign.Name)
				}
				if !strings.Contains(out, strconv.Itoa(resIDAlign.ID)) {
					t.Errorf("%d, not found", resIDAlign.ID)
				}
			}
			fmt.Println(out)
		})
	}
}

func TestEmptyTables(t *testing.T) {
	emptyResourcesList := []*meroxa.Resource{}
	out := CaptureOutput(func() {
		PrintResourcesTable(emptyResourcesList)
	})

	if out != "" {
		t.Errorf("Output for resources should be blank")
	}

	emptyConnectorsList := []*meroxa.Connector{}
	out = CaptureOutput(func() {
		PrintConnectorsTable(emptyConnectorsList)
	})

	if out != "" {
		t.Errorf("Output for connectors should be blank")
	}

	emptyPipelinesList := []*meroxa.Pipeline{}

	out = CaptureOutput(func() {
		PrintPipelinesTable(emptyPipelinesList)
	})

	if out != "" {
		t.Errorf("Output for pipelines should be blank")
	}
}
func TestResourceTypesTable(t *testing.T) {
	types := []string{"postgres", "s3", "redshift", "mysql", "jdbc", "url", "mongodb"}
	PrintResourceTypesTable(types)
}

func TestConnectionsTable(t *testing.T) {
	connectionIDAlign := &meroxa.Connector{}
	connectionInputOutput := &meroxa.Connector{}
	connection := &meroxa.Connector{
		ID:            0,
		Type:          "jdbc",
		Name:          "base",
		Configuration: nil,
		Metadata:      nil,
		Streams: map[string]interface{}{
			"dynamic": "false",
			"output":  []interface{}{"output-foo", "output-bar"},
		},
		State:      "running",
		Trace:      "",
		PipelineID: 1,
	}

	deepCopy(connection, connectionIDAlign)
	connectionIDAlign.Name = "id-alignment"
	connectionIDAlign.ID = 1000

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

	for name, connections := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintConnectorsTable(connections)
			})

			switch name {
			case "Base":
				if !strings.Contains(out, connection.Name) {
					t.Errorf("%s, not found", connection.Name)
				}
				if !strings.Contains(out, strconv.Itoa(connection.ID)) {
					t.Errorf("%d, not found", connection.ID)
				}
			case "ID_Alignment":
				if !strings.Contains(out, connectionIDAlign.Name) {
					t.Errorf("%s, not found", connectionIDAlign.Name)
				}
				if !strings.Contains(out, strconv.Itoa(connectionIDAlign.ID)) {
					t.Errorf("%d, not found", connectionIDAlign.ID)
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

func TestPipelinesTable(t *testing.T) {
	pipelineIDAlign := &meroxa.Pipeline{}
	pipelineBase := &meroxa.Pipeline{
		ID:       0,
		Name:     "pipeline-base",
		Metadata: nil,
	}
	deepCopy(pipelineBase, pipelineIDAlign)
	pipelineIDAlign.ID = 1000
	pipelineIDAlign.Name = "pipeline-align"

	tests := map[string][]*meroxa.Pipeline{
		"Base":         {pipelineBase},
		"ID_Alignment": {pipelineBase, pipelineIDAlign},
	}

	for name, pipelines := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintPipelinesTable(pipelines)
			})
			switch name {
			case "Base":
				if !strings.Contains(out, pipelineBase.Name) {
					t.Errorf("%s, not found", pipelineBase.Name)
				}
				if !strings.Contains(out, strconv.Itoa(pipelineBase.ID)) {
					t.Errorf("%d, not found", pipelineBase.ID)
				}
			case "ID_Alignment":
				if !strings.Contains(out, pipelineBase.Name) {
					t.Errorf("%s, not found", pipelineBase.Name)
				}
				if !strings.Contains(out, strconv.Itoa(pipelineBase.ID)) {
					t.Errorf("%d, not found", pipelineBase.ID)
				}
			}
			fmt.Println(out)
		})
	}
}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	_ = json.Unmarshal(byt, b)
}
