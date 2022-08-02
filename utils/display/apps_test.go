package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/volatiletech/null/v8"

	"github.com/meroxa/cli/utils"
)

func TestAppLogsTable(t *testing.T) {
	app := utils.GenerateApplication("")
	res1 := "res1"
	res2 := "res2"

	app.Resources = []meroxa.ApplicationResource{
		{
			EntityIdentifier: meroxa.EntityIdentifier{Name: null.StringFrom(res1)},
			Collection: meroxa.ResourceCollection{
				Name:   null.StringFrom(res1),
				Source: null.StringFrom("source"),
			},
		},
		{
			EntityIdentifier: meroxa.EntityIdentifier{Name: null.StringFrom(res2)},
			Collection: meroxa.ResourceCollection{
				Name:        null.StringFrom(res2),
				Destination: null.StringFrom("destination"),
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

	out := AppLogsTable(app.Resources, connectors, functions)

	if !strings.Contains(out, "custom log") {
		t.Errorf("expected %q to be shown", log)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (source)", res1)) {
		t.Errorf("expected %q to be shown as a source", res1)
	}

	if !strings.Contains(out, fmt.Sprintf("%s (destination)", res2)) {
		t.Errorf("expected %q to be shown as a destination", res2)
	}
}
