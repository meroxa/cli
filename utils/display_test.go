package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestResourcesTable(t *testing.T) {
	resource := &meroxa.Resource{
		ID:          0,
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "error",
		},
	}
	resIDAlign := &meroxa.Resource{
		ID:          1000,
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "ready",
		},
	}

	tests := map[string][]*meroxa.Resource{
		"Base":         {resource},
		"ID_Alignment": {resource, resIDAlign},
	}

	tableHeaders := []string{"ID", "NAME", "TYPE", "ENVIRONMENT", "URL", "TUNNEL", "STATE"}

	for name, resources := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintResourcesTable(resources, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "Base":
				if !strings.Contains(out, resource.Name) {
					t.Errorf("%s, not found", resource.Name)
				}
				if !strings.Contains(out, strconv.Itoa(resource.ID)) {
					t.Errorf("%d, not found", resource.ID)
				}
				if !strings.Contains(out, strings.Title(string(resource.Status.State))) {
					t.Errorf("state %s, not found", resource.Status.State)
				}
			case "ID_Alignment":
				if !strings.Contains(out, resIDAlign.Name) {
					t.Errorf("%s, not found", resIDAlign.Name)
				}
				if !strings.Contains(out, strconv.Itoa(resIDAlign.ID)) {
					t.Errorf("%d, not found", resIDAlign.ID)
				}
				if !strings.Contains(out, strings.Title(string(resIDAlign.Status.State))) {
					t.Errorf("state %s, not found", resource.Status.State)
				}
			}
			fmt.Println(out)
		})
	}
}

func TestResourcesTableWithoutHeaders(t *testing.T) {
	resource := &meroxa.Resource{
		ID:          0,
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "error",
		},
	}

	var resources []*meroxa.Resource
	resources = append(resources, resource)

	tableHeaders := []string{"ID", "NAME", "TYPE", "URL", "TUNNEL", "STATE"}

	out := CaptureOutput(func() {
		PrintResourcesTable(resources, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, resource.Name) {
		t.Errorf("%s, not found", resource.Name)
	}
	if !strings.Contains(out, strconv.Itoa(resource.ID)) {
		t.Errorf("%d, not found", resource.ID)
	}
	if !strings.Contains(out, strings.Title(string(resource.Status.State))) {
		t.Errorf("state %s, not found", resource.Status.State)
	}
}

func TestEmptyTables(t *testing.T) {
	var emptyResourcesList []*meroxa.Resource
	out := CaptureOutput(func() {
		PrintResourcesTable(emptyResourcesList, true)
	})

	if out != "\n" {
		t.Errorf("Output for resources should be blank")
	}

	var emptyConnectorsList []*meroxa.Connector
	out = CaptureOutput(func() {
		PrintConnectorsTable(emptyConnectorsList, true)
	})

	if out != "\n" {
		t.Errorf("Output for connectors should be blank")
	}

	var emptyPipelinesList []*meroxa.Pipeline
	out = CaptureOutput(func() {
		PrintPipelinesTable(emptyPipelinesList, true)
	})

	if out != "\n" {
		t.Errorf("Output for pipelines should be blank")
	}

	var emptyEnvironmentsList []*meroxa.Environment
	out = CaptureOutput(func() {
		PrintEnvironmentsTable(emptyEnvironmentsList, true)
	})

	if out != "\n" {
		t.Errorf("Output for pipelines should be blank")
	}
}
func TestResourceTypesTable(t *testing.T) {
	types := []string{"postgres", "s3", "redshift", "mysql", "jdbc", "url", "mongodb"}

	out := CaptureOutput(func() {
		PrintResourceTypesTable(types, false)
	})

	if !strings.Contains(out, "TYPES") {
		t.Errorf("table headers is missing")
	}

	for _, rType := range types {
		if !strings.Contains(out, rType) {
			t.Errorf("%s, not found", rType)
		}
	}
}

func TestResourceTypesTableWithoutHeaders(t *testing.T) {
	types := []string{"postgres", "s3", "redshift", "mysql", "jdbc", "url", "mongodb"}
	out := CaptureOutput(func() {
		PrintResourceTypesTable(types, true)
	})

	if strings.Contains(out, "TYPES") {
		t.Errorf("table header should not be displayed")
	}

	for _, rType := range types {
		if !strings.Contains(out, rType) {
			t.Errorf("%s, not found", rType)
		}
	}
}

func TestConnectorRunningTable(t *testing.T) {
	connector := &meroxa.Connector{
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
	failedConnector := &meroxa.Connector{}
	deepCopy(connector, failedConnector)
	failedConnector.State = "failed"
	failedConnector.Trace = "exception goes here"

	tests := map[string]*meroxa.Connector{
		"running": connector,
		"failed":  failedConnector,
	}

	tableHeaders := []string{"ID", "Name", "Type", "Streams", "State", "Pipeline"}

	for name, connector := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				fmt.Println(ConnectorTable(connector))
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "running":
				if !strings.Contains(out, connector.Name) {
					t.Errorf("%s, not found", connector.Name)
				}
				if !strings.Contains(out, strconv.Itoa(connector.ID)) {
					t.Errorf("%d, not found", connector.ID)
				}
			case "failed":
				if !strings.Contains(out, connector.Name) {
					t.Errorf("%s, not found", connector.Name)
				}
				if !strings.Contains(out, strconv.Itoa(connector.ID)) {
					t.Errorf("%d, not found", connector.ID)
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

	tableHeaders := []string{"ID", "NAME", "TYPE", "STREAMS", "STATE", "PIPELINE"}

	for name, connections := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintConnectorsTable(connections, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

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

func TestConnectorsTableWithoutHeaders(t *testing.T) {
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

	tableHeaders := []string{"ID", "NAME", "TYPE", "STREAMS", "STATE", "PIPELINE"}

	var connections []*meroxa.Connector
	connections = append(connections, connection)

	out := CaptureOutput(func() {
		PrintConnectorsTable(connections, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, connection.Name) {
		t.Errorf("%s, not found", connection.Name)
	}
	if !strings.Contains(out, strconv.Itoa(connection.ID)) {
		t.Errorf("%d, not found", connection.ID)
	}
}

func TestPipelinesTable(t *testing.T) {
	pipelineIDAlign := &meroxa.Pipeline{}
	pipelineWithEnv := &meroxa.Pipeline{}

	pipelineBase := &meroxa.Pipeline{
		UUID: "6f380820-dfed-4a69-b708-10d134866a35",
		ID:   0,
		Name: "pipeline-base",
	}
	deepCopy(pipelineBase, pipelineIDAlign)
	pipelineIDAlign.UUID = "0e1d29b9-2e62-4cc2-a49d-126f2e1b15ef"
	pipelineIDAlign.Name = "pipeline-align"
	pipelineIDAlign.ID = 1000

	deepCopy(pipelineBase, pipelineWithEnv)
	pipelineWithEnv.UUID = "038de172-c4b0-49d8-a1d9-26fbeaa2f726"
	pipelineWithEnv.Environment = &meroxa.PipelineEnvironment{
		UUID: "e56b1b2e-b6d7-455d-887e-84a0823d84a8",
		Name: "my-environment",
	}

	tests := map[string][]*meroxa.Pipeline{
		"Base":             {pipelineBase},
		"ID_Alignment":     {pipelineBase, pipelineIDAlign},
		"With_Environment": {pipelineBase, pipelineIDAlign, pipelineWithEnv},
	}

	tableHeaders := []string{"UUID", "ID", "NAME", "ENVIRONMENT", "STATE"}

	for name, pipelines := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintPipelinesTable(pipelines, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "Base":
				if !strings.Contains(out, pipelineBase.Name) {
					t.Errorf("%s, not found", pipelineBase.Name)
				}
				if !strings.Contains(out, strconv.Itoa(pipelineBase.ID)) {
					t.Errorf("%d, not found", pipelineBase.ID)
				}
				if !strings.Contains(out, string(meroxa.EnvironmentTypeCommon)) {
					t.Errorf("environment should be %s", string(meroxa.EnvironmentTypeCommon))
				}
			case "ID_Alignment":
				if !strings.Contains(out, pipelineIDAlign.Name) {
					t.Errorf("%s, not found", pipelineIDAlign.Name)
				}
				if !strings.Contains(out, strconv.Itoa(pipelineIDAlign.ID)) {
					t.Errorf("%d, not found", pipelineIDAlign.ID)
				}
			case "With_Environment":
				if !strings.Contains(out, pipelineWithEnv.Environment.Name) {
					t.Errorf("expected environment name to be %q", pipelineWithEnv.Environment.Name)
				}
			}

			fmt.Println(out)
		})
	}
}

func TestPipelinesTableWithoutHeaders(t *testing.T) {
	pipeline := &meroxa.Pipeline{
		ID:   0,
		Name: "pipeline-base",
	}

	var pipelines []*meroxa.Pipeline
	tableHeaders := []string{"ID", "NAME", "STATE"}

	pipelines = append(pipelines, pipeline)

	out := CaptureOutput(func() {
		PrintPipelinesTable(pipelines, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, pipeline.Name) {
		t.Errorf("%s, not found", pipeline.Name)
	}
	if !strings.Contains(out, strconv.Itoa(pipeline.ID)) {
		t.Errorf("%d, not found", pipeline.ID)
	}
}

func TestEnvironmentsTable(t *testing.T) {
	e := &meroxa.Environment{
		Type:     meroxa.EnvironmentTypeDedicated,
		Name:     "environment-1234",
		Provider: meroxa.EnvironmentProviderAws,
		Region:   meroxa.EnvironmentRegionUsEast2,
		Status:   meroxa.EnvironmentViewStatus{State: meroxa.EnvironmentStateProvisioned},
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
	}

	tests := map[string][]*meroxa.Environment{
		"Base": {e},
	}

	tableHeaders := []string{"ID", "NAME", "TYPE", "PROVIDER", "REGION", "STATE"}

	for name, environments := range tests {
		t.Run(name, func(t *testing.T) {
			out := CaptureOutput(func() {
				PrintEnvironmentsTable(environments, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}
			if !strings.Contains(out, e.Name) {
				t.Errorf("%s, not found", e.Name)
			}
			if !strings.Contains(out, string(e.Type)) {
				t.Errorf("%s, not found", e.Type)
			}
			if !strings.Contains(out, string(e.Region)) {
				t.Errorf("%s, not found", e.Region)
			}
			if !strings.Contains(out, string(e.Status.State)) {
				t.Errorf("%s, not found", e.Status.State)
			}
			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}

			fmt.Println(out)
		})
	}
}

func TestEnvironmentsTableWithoutHeaders(t *testing.T) {
	e := &meroxa.Environment{
		Type:     meroxa.EnvironmentTypeDedicated,
		Name:     "environment-1234",
		Provider: meroxa.EnvironmentProviderAws,
		Region:   meroxa.EnvironmentRegionUsEast2,
		Status:   meroxa.EnvironmentViewStatus{State: meroxa.EnvironmentStateProvisioned},
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
	}

	var environments []*meroxa.Environment
	tableHeaders := []string{"ID", "NAME", "TYPE", "PROVIDER", "REGION", "STATE"}

	environments = append(environments, e)

	out := CaptureOutput(func() {
		PrintEnvironmentsTable(environments, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, e.UUID) {
		t.Errorf("%s, not found", e.UUID)
	}
	if !strings.Contains(out, e.Name) {
		t.Errorf("%s, not found", e.Name)
	}
	if !strings.Contains(out, string(e.Type)) {
		t.Errorf("%s, not found", e.Type)
	}
	if !strings.Contains(out, string(e.Region)) {
		t.Errorf("%s, not found", e.Region)
	}
	if !strings.Contains(out, string(e.Status.State)) {
		t.Errorf("%s, not found", e.Status.State)
	}
	if !strings.Contains(out, e.UUID) {
		t.Errorf("%s, not found", e.UUID)
	}
}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	_ = json.Unmarshal(byt, b)
}
