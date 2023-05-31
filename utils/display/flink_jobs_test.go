package display

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/stretchr/testify/assert"
)

func TestFlinkJobsTable(t *testing.T) {
	fOne := &meroxa.FlinkJob{
		UUID:          "424ec647-9f0f-45a5-8e4b-3e0441f12555",
		Name:          "my-flink-job",
		InputStreams:  []string{"inputstream_one", "inputstream_two"},
		OutputStreams: []string{"outtt"},
		Status: meroxa.FlinkJobStatus{
			State:                  "running",
			LifecycleState:         "success",
			ReconciliationState:    "deployed",
			ManagerDeploymentState: "ready",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	fTwo := &meroxa.FlinkJob{
		UUID:          "123d4da3-9f0f-45a5-8e4b-77777777",
		Name:          "squirrel-app",
		InputStreams:  []string{"inputstream_one"},
		OutputStreams: []string{"outtt", "anotheroutt"},
		Status: meroxa.FlinkJobStatus{
			State:                  "failed",
			LifecycleState:         "suspended",
			ReconciliationState:    "rolling back",
			ManagerDeploymentState: "error",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	tests := map[string][]*meroxa.FlinkJob{
		"Base":         {fOne},
		"ID_Alignment": {fOne, fTwo},
		"Empty":        {},
	}

	for name, flinkJobs := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintFlinkJobsTable(flinkJobs)
			})

			switch name {
			case "Base":
				verifyPrintFlinkJobsOutput(t, out, fOne)
			case "ID_Alignment":
				verifyPrintFlinkJobsOutput(t, out, fOne)
				verifyPrintFlinkJobsOutput(t, out, fTwo)
			case "Empty":
				assert.Equal(t, out, "\n")
			}
			fmt.Println(out)
		})
	}
}

func verifyPrintFlinkJobsOutput(t *testing.T, out string, flinkJob *meroxa.FlinkJob) {
	// verify header fields
	tableHeaders := []string{"UUID", "NAME", "STATE", "DEPLOYMENT STATE"}

	for _, header := range tableHeaders {
		if !strings.Contains(out, header) {
			t.Errorf("%s header is missing", header)
		}
	}

	// verify fields that are supposed to be included in the output
	if !strings.Contains(out, flinkJob.Name) {
		t.Errorf("%s, not found", flinkJob.Name)
	}
	if !strings.Contains(out, flinkJob.UUID) {
		t.Errorf("%s, not found", flinkJob.UUID)
	}
	if !strings.Contains(out, string(flinkJob.Status.State)) {
		t.Errorf("state %s, not found", flinkJob.Status.State)
	}
	// verify fields that are supposed to be excluded from the output
	if strings.Contains(out, fmt.Sprintf("%v", flinkJob.InputStreams)) {
		t.Errorf("found unwanted output: %s", flinkJob.InputStreams)
	}
	if strings.Contains(out, fmt.Sprintf("%v", flinkJob.OutputStreams)) {
		t.Errorf("found unwanted output: %s", flinkJob.OutputStreams)
	}
	if strings.Contains(out, string(flinkJob.Status.LifecycleState)) {
		t.Errorf("found unwanted output: %s", string(flinkJob.Status.LifecycleState))
	}
	if strings.Contains(out, string(flinkJob.Status.ReconciliationState)) {
		t.Errorf("found unwanted output: %s", string(flinkJob.Status.ReconciliationState))
	}
	if strings.Contains(out, flinkJob.CreatedAt.String()) {
		t.Errorf("found unwanted output: %s", flinkJob.CreatedAt.String())
	}
	if strings.Contains(out, flinkJob.UpdatedAt.String()) {
		t.Errorf("found unwanted output: %s", flinkJob.UpdatedAt.String())
	}
}
