package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestPipelinesTable(t *testing.T) {
	pipelineIDAlign := &meroxa.Pipeline{}
	pipelineWithEnv := &meroxa.Pipeline{}

	pipelineBase := &meroxa.Pipeline{
		UUID: "6f380820-dfed-4a69-b708-10d134866a35",
		Name: "pipeline-base",
	}
	deepCopy(pipelineBase, pipelineIDAlign)
	pipelineIDAlign.UUID = "0e1d29b9-2e62-4cc2-a49d-126f2e1b15ef"
	pipelineIDAlign.Name = "pipeline-align"

	deepCopy(pipelineBase, pipelineWithEnv)
	pipelineWithEnv.UUID = "038de172-c4b0-49d8-a1d9-26fbeaa2f726"
	pipelineWithEnv.Environment = &meroxa.EntityIdentifier{
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
			out := utils.CaptureOutput(func() {
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
				if !strings.Contains(out, pipelineBase.UUID) {
					t.Errorf("%s, not found", pipelineBase.UUID)
				}
				if !strings.Contains(out, string(meroxa.EnvironmentTypeCommon)) {
					t.Errorf("environment should be %s", string(meroxa.EnvironmentTypeCommon))
				}
			case "ID_Alignment":
				if !strings.Contains(out, pipelineIDAlign.Name) {
					t.Errorf("%s, not found", pipelineIDAlign.Name)
				}
				if !strings.Contains(out, pipelineIDAlign.UUID) {
					t.Errorf("%s, not found", pipelineIDAlign.UUID)
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

func TestPipelineTable(t *testing.T) {
	pipelineWithEnv := &meroxa.Pipeline{}

	pipelineBase := &meroxa.Pipeline{
		UUID: "6f380820-dfed-4a69-b708-10d134866a35",
		Name: "pipeline-base",
	}

	deepCopy(pipelineBase, pipelineWithEnv)
	pipelineWithEnv.UUID = "038de172-c4b0-49d8-a1d9-26fbeaa2f726"
	pipelineWithEnv.Environment = &meroxa.EntityIdentifier{
		UUID: "e56b1b2e-b6d7-455d-887e-84a0823d84a8",
		Name: "my-environment",
	}

	tests := map[string]*meroxa.Pipeline{
		"Base":             pipelineBase,
		"With_Environment": pipelineWithEnv,
	}

	tableHeaders := []string{"UUID", "ID", "Name", "State"}
	envHeader := "Environment Name"

	for name, p := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintPipelineTable(p)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%q header is missing", header)
				}
			}

			switch name {
			case "Base":
				if !strings.Contains(out, pipelineBase.Name) {
					t.Errorf("%s, not found", pipelineBase.Name)
				}
				if !strings.Contains(out, pipelineBase.UUID) {
					t.Errorf("%s, not found", pipelineBase.UUID)
				}
				if !strings.Contains(out, pipelineBase.UUID) {
					t.Errorf("%s, not found", pipelineBase.UUID)
				}
				if !strings.Contains(out, envHeader) {
					t.Errorf("%q not found", envHeader)
				}
			case "With_Environment":
				if !strings.Contains(out, pipelineWithEnv.Environment.UUID) {
					t.Errorf("expected environment UUID to be %q", pipelineWithEnv.Environment.UUID)
				}
			}
			fmt.Println(out)
		})
	}
}

func TestPipelinesTableWithoutHeaders(t *testing.T) {
	pipeline := &meroxa.Pipeline{
		UUID: "6f380820-dfed-4a69-b708-10d134866a35",
		Name: "pipeline-base",
	}

	var pipelines []*meroxa.Pipeline
	tableHeaders := []string{"ID", "NAME", "STATE"}

	pipelines = append(pipelines, pipeline)

	out := utils.CaptureOutput(func() {
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
	if !strings.Contains(out, pipeline.UUID) {
		t.Errorf("%s, not found", pipeline.UUID)
	}
}
