package display

import (
	"testing"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestEmptyTables(t *testing.T) {
	var emptyResourcesList []*meroxa.Resource
	out := utils.CaptureOutput(func() {
		PrintResourcesTable(emptyResourcesList, true)
	})

	if out != "\n" {
		t.Errorf("Output for resources should be blank")
	}

	var emptyConnectorsList []*meroxa.Connector
	out = utils.CaptureOutput(func() {
		PrintConnectorsTable(emptyConnectorsList, true)
	})

	if out != "\n" {
		t.Errorf("Output for connectors should be blank")
	}

	var emptyPipelinesList []*meroxa.Pipeline
	out = utils.CaptureOutput(func() {
		PrintPipelinesTable(emptyPipelinesList, true)
	})

	if out != "\n" {
		t.Errorf("Output for pipelines should be blank")
	}

	var emptyEnvironmentsList []*meroxa.Environment
	out = utils.CaptureOutput(func() {
		PrintEnvironmentsTable(emptyEnvironmentsList, true)
	})

	if out != "\n" {
		t.Errorf("Output for pipelines should be blank")
	}
}
