package cmd

import (
	"encoding/json"
	"github.com/meroxa/meroxa-go"
	"testing"
)

func TestResourcesTable(t *testing.T) {
	resource := &meroxa.Resource{
		ID:            0,
		Kind:          "jdbc",
		Name:          "my-db-jdbc-source",
		URL:           "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials:   nil,
		Configuration: nil,
		Metadata:      nil,
	}
	resIDLong := &meroxa.Resource{
		ID:            1000,
		Kind:          "jdbc",
		Name:          "my-db-jdbc-source",
		URL:           "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials:   nil,
		Configuration: nil,
		Metadata:      nil,
	}

	tests := map[string][]*meroxa.Resource{
		"Base":         {resource},
		"ID Alignment": {resource, resIDLong},
	}

	for name, resources := range tests {
		t.Run(name, func(t *testing.T) {
			printResourcesTable(resources)
		})
	}
}

func TestResourceTypesTable(t *testing.T) {
	types := []string{"postgres", "s3", "redshift", "mysql", "jdbc", "url", "mongodb"}
	printResourceTypesTable(types)
}

func TestConnectionsTable(t *testing.T) {
	connectionIDAlign := &meroxa.Connector{}
	connectionInputOutput := &meroxa.Connector{}
	connection := &meroxa.Connector{
		ID:            0,
		Kind:          "jdbc",
		Name:          "base",
		Configuration: nil,
		Metadata:      nil,
		Streams: map[string]interface{}{
			"dynamic": "false",
			"output":  []interface{}{"output-foo", "output-bar"},
		},
		State: "running",
		Trace: "",
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
			printConnectionsTable(connections)
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

	for name, pipeline := range tests {
		t.Run(name, func(t *testing.T) {
			printPipelinesTable(pipeline)
		})
	}
}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}
