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
	tableHeaders := []string{"UUID", "NAME", "LIFECYCLE STATE"}

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
	if !strings.Contains(out, string(flinkJob.Status.LifecycleState)) {
		t.Errorf("state %s, not found", flinkJob.Status.LifecycleState)
	}
	if !strings.Contains(out, string(meroxa.EnvironmentTypeCommon)) {
		t.Errorf("state %s, not found", string(meroxa.EnvironmentTypeCommon))
	}
	// verify fields that are supposed to be excluded from the output
	if strings.Contains(out, fmt.Sprintf("%v", flinkJob.InputStreams)) {
		t.Errorf("found unwanted output: %s", flinkJob.InputStreams)
	}
	if strings.Contains(out, fmt.Sprintf("%v", flinkJob.OutputStreams)) {
		t.Errorf("found unwanted output: %s", flinkJob.OutputStreams)
	}
	if strings.Contains(out, string(flinkJob.Status.State)) {
		t.Errorf("found unwanted output: %s", string(flinkJob.Status.State))
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

func TestFlinkJobTable(t *testing.T) {
	testCases := []struct {
		desc                 string
		flinkJob             func() *meroxa.FlinkJob
		shouldErrorOnEnvInfo func(string) bool
		checkStates          bool
	}{
		{
			desc: "Flink Job with no environment",
			flinkJob: func() *meroxa.FlinkJob {
				a := utils.GenerateFlinkJob()
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, string(meroxa.EnvironmentTypeCommon))
			},
		},
		{
			desc: "Flink Job with in a private environment",
			flinkJob: func() *meroxa.FlinkJob {
				a := utils.GenerateFlinkJob()
				a.Environment.Name = "hey-now"
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return strings.Contains(output, string(meroxa.EnvironmentTypeCommon))
			},
		},
		{
			desc: "Flink Job with states and details",
			flinkJob: func() *meroxa.FlinkJob {
				a := utils.GenerateFlinkJob()
				a.Status.State = meroxa.FlinkJobStateRunning
				a.Status.ManagerDeploymentState = meroxa.FlinkJobManagerDeploymentStateDeploying
				a.Status.LifecycleState = meroxa.FlinkJobLifecycleStateCreated
				a.Status.ReconciliationState = meroxa.FlinkJobReconciliationStateDeployed
				a.Status.Details = "so many good things"
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, string(meroxa.EnvironmentTypeCommon))
			},
			checkStates: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fj := tc.flinkJob()
			out := FlinkJobTable(fj)
			if !strings.Contains(out, "UUID") {
				t.Errorf("expected %q to be shown\n%s\n", fj.UUID, out)
			}
			if !strings.Contains(out, "Name:") {
				t.Errorf("expected %q to be shown\n%s\n", fj.Name, out)
			}
			if !strings.Contains(out, "Created At:") {
				t.Errorf("expected %q to be shown\n%s\n", fj.CreatedAt, out)
			}
			if !strings.Contains(out, "Updated At:") {
				t.Errorf("expected %q to be shown\n%s\n", fj.UpdatedAt, out)
			}
			if !strings.Contains(out, "Input Streams:") {
				t.Errorf("expected %q to be shown\n%s\n", fj.InputStreams, out)
			}
			if !strings.Contains(out, "Output Streams:") {
				t.Errorf("expected %q to be shown\n%s\n", fj.OutputStreams, out)
			}
			if !strings.Contains(out, "Lifecycle State") {
				t.Errorf("expected %q to be shown\n%s\n", fj.Status.LifecycleState, out)
			}
			if tc.checkStates {
				if !strings.Contains(out, "Job State:") {
					t.Errorf("expected %q to be shown\n%s\n", fj.Status.State, out)
				}
				if !strings.Contains(out, "Reconciliation State:") {
					t.Errorf("expected %q to be shown\n%s\n", fj.Status.ReconciliationState, out)
				}
				if !strings.Contains(out, "Manager Deployment State:") {
					t.Errorf("expected %q to be shown\n%s\n", fj.Status.ManagerDeploymentState, out)
				}
				if !strings.Contains(out, "State Details:") {
					t.Errorf("expected %q to be shown\n%s\n", fj.Status.Details, out)
				}
			}
			if tc.shouldErrorOnEnvInfo(out) {
				t.Errorf("expected environment information to be shown\n%s\n", out)
			}
		})
	}
}
