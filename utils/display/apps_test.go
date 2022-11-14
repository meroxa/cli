package display

import (
	"fmt"
	"strings"
	"testing"

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
