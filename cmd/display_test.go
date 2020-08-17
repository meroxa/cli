package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meroxa/meroxa-go"
	"io"
	"os"
	"strconv"
	"strings"
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
	resIDAlign := &meroxa.Resource{
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
		"ID_Alignment": {resource, resIDAlign},
	}

	for name, resources := range tests {
		t.Run(name, func(t *testing.T) {
			out := captureOutput(func() {
				printResourcesTable(resources)
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
			out := captureOutput(func() {
				printConnectionsTable(connections)
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
			out := captureOutput(func() {
				printPipelinesTable(pipelines)
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
	json.Unmarshal(byt, b)
}

func captureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	print()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	f()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out
}
