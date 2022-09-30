package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/meroxa/cli/utils"
)

func TestAppLogsTable(t *testing.T) {
	app := utils.GenerateApplication("")
	res1 := "res1"
	res2 := "res2"

	app.Resources = []meroxa.ApplicationResource{
		{
			EntityIdentifier: meroxa.EntityIdentifier{Name: res1},
			Collection: meroxa.ResourceCollection{
				Name:   res1,
				Source: "source",
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{Name: res2},
			Collection: meroxa.ResourceCollection{
				Name:        res2,
				Destination: "destination",
			},
		},
	}

	log := "custom log"

	connectors := []*AppExtendedConnector{
		{Connector: &meroxa.Connector{
			Name: "conn1", ResourceName: res1, Type: meroxa.ConnectorTypeSource, State: meroxa.ConnectorStateRunning},
			Logs: log},
		{Connector: &meroxa.Connector{
			Name: "conn2", ResourceName: res2, Type: meroxa.ConnectorTypeDestination, State: meroxa.ConnectorStateRunning},
			Logs: log},
	}

	functions := []*meroxa.Function{
		{Name: "fun1", UUID: "abc-def", Status: meroxa.FunctionStatus{State: "running"}, Logs: log},
	}

	deployment := &meroxa.Deployment{
		UUID:   "ghi-jkl",
		Status: meroxa.DeploymentStatus{Details: "deployment in progress"}}

	out := AppLogsTable(app.Resources, connectors, functions, deployment)

	if !strings.Contains(out, "custom log") {
		t.Errorf("expected %q to be shown", log)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (source)", res1)) {
		t.Errorf("expected %q to be shown as a source", res1)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (destination)", res2)) {
		t.Errorf("expected %q to be shown as a destination", res2)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (function)", functions[0].Name)) {
		t.Errorf("expected %q to be shown as a function", functions[0].Name)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (deployment)", deployment.UUID)) {
		t.Errorf("expected %q to be shown as a deployment", deployment.UUID)
	}
}
