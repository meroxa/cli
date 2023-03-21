package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestAppLogsTable(t *testing.T) {
	res1 := "res1"
	res2 := "res2"
	fun1 := "fun1"
	deploymentUUID := "uu-id"
	log := "custom log"
	logs := meroxa.ApplicationLogs{
		ConnectorLogs:  map[string]string{"source " + res1: log, "destination " + res2: log},
		FunctionLogs:   map[string]string{fun1: log},
		DeploymentLogs: map[string]string{deploymentUUID: log},
	}

	out := AppLogsTable(&logs)

	if !strings.Contains(out, "custom log") {
		t.Errorf("expected %q to be shown", log)
	}

	if !strings.Contains(out, fmt.Sprintf("# Logs for source %s resource\n\n%s", res1, log)) {
		t.Errorf("expected %q to be shown with logs", res1)
	}
	if !strings.Contains(out, fmt.Sprintf("# Logs for destination %s resource\n\n%s", res2, log)) {
		t.Errorf("expected %q to be shown with logs", res2)
	}
	if !strings.Contains(out, fmt.Sprintf("# Logs for %s function\n\n%s", fun1, log)) {
		t.Errorf("expected %q to be shown with logs", fun1)
	}
	if !strings.Contains(out, fmt.Sprintf("# Logs for %s deployment\n\n%s", deploymentUUID, log)) {
		t.Errorf("expected %q to be shown with logs", deploymentUUID)
	}
}

func TestAppDescribeTable(t *testing.T) {
	testCases := []struct {
		desc                 string
		app                  func() *meroxa.Application
		shouldErrorOnEnvInfo func(string) bool
	}{
		{
			desc: "Application with no environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplication("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return strings.Contains(output, "Environment")
			},
		},
		{
			desc: "Application with in a private environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplicationWithEnv("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, "Environment")
			},
		},
		{
			desc: "Application with in a private environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplicationWithEnv("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, "Environment")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			app := tc.app()
			out := AppTable(app)
			if !strings.Contains(out, "UUID") {
				t.Errorf("expected %q to be shown\n%s\n", app.UUID, out)
			}
			if !strings.Contains(out, "Name:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Name, out)
			}
			if !strings.Contains(out, "Language:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Language, out)
			}
			if !strings.Contains(out, "Git SHA:") {
				t.Errorf("expected %q to be shown\n%s\n", app.GitSha, out)
			}
			if !strings.Contains(out, "Created At:") {
				t.Errorf("expected %q to be shown\n%s\n", app.CreatedAt, out)
			}
			if !strings.Contains(out, "Updated At:") {
				t.Errorf("expected %q to be shown\n%s\n", app.UpdatedAt, out)
			}
			if !strings.Contains(out, "State:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Status.State, out)
			}
			if tc.shouldErrorOnEnvInfo(out) {
				t.Errorf("expected environment information to be shown\n%s\n", out)
			}
		})
	}
}
