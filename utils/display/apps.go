package display

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type AppExtendedConnector struct {
	Connector *meroxa.Connector
	Logs      string
}

func AppsTable(apps []*meroxa.Application, hideHeaders bool) string {
	if len(apps) == 0 {
		return ""
	}

	table := simpletable.New()
	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "LANGUAGE"},
				{Align: simpletable.AlignCenter, Text: "GIT SHA"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
			},
		}
	}

	for _, app := range apps {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: app.UUID},
			{Align: simpletable.AlignCenter, Text: app.Name},
			{Align: simpletable.AlignCenter, Text: app.Language},
			{Align: simpletable.AlignCenter, Text: app.GitSha},
			{Align: simpletable.AlignCenter, Text: string(app.Status.State)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func AppTable(app *meroxa.Application) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: app.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: app.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Language:"},
			{Text: app.Language},
		},
		{
			{Align: simpletable.AlignRight, Text: "Git SHA:"},
			{Text: strings.TrimSpace(app.GitSha)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: app.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: app.UpdatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(app.Status.State)},
		},
	}
	mainTable.SetStyle(simpletable.StyleCompact)
	output := mainTable.String()

	subTable := appResourcesTable(app.Resources, app.Connectors)
	if subTable != "" {
		output += "\n" + subTable
	}
	subTable = appFunctionsTable(app.Functions)
	if subTable != "" {
		output += "\n" + subTable
	}
	return output
}

func isMatching(collection meroxa.ResourceCollection, connectorType string) bool {
	t := "destination"
	if collection.Source == "true" {
		t = "source"
	}
	return strings.Contains(connectorType, t)
}

func appResourcesTable(resources []meroxa.ApplicationResource, connectors []meroxa.EntityDetails) string {
	if len(resources) == 0 {
		return ""
	}
	subTable := "\tResources\n"
	var connector *meroxa.EntityDetails

	for _, r := range resources {
		t := "source"
		if r.Collection.Destination == "true" {
			t = "destination"
		}
		for _, c := range connectors {
			if r.UUID == c.ResourceUUID && isMatching(r.Collection, c.ResourceType) {
				connector = &c
				break
			}
		}

		subTable += fmt.Sprintf("\t    %s (%s)\n", r.Name, t)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", r.UUID)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "Type", r.ResourceType)
		if connector != nil {
			subTable += fmt.Sprintf("\t\t%5s:   %s\n", "State", connector.Status)
		}
		connector = nil
	}

	return subTable
}

func appFunctionsTable(functions []meroxa.EntityDetails) string {
	if len(functions) == 0 {
		return ""
	}
	subTable := "\tFunctions\n"

	for _, f := range functions {
		subTable += fmt.Sprintf("\t    %s\n", f.Name)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", f.UUID)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "State", f.Status)
	}

	return subTable
}

func AppLogsTable(
	resources []meroxa.ApplicationResource,
	connectors []*AppExtendedConnector,
	functions []*meroxa.Function,
	deployment *meroxa.Deployment) string {
	var r meroxa.ApplicationResource
	var subTable string

	for _, c := range connectors {
		found := false
		for _, resource := range resources {
			if resource.Name == c.Connector.ResourceName {
				r = resource
				found = true
				break
			}
		}
		if !found {
			panic("internal error")
		}

		// Only show information if there are logs or a trace available
		if c.Logs != "" || c.Connector.Trace != "" {
			if r.Collection.Source != "" {
				subTable += fmt.Sprintf("\n# Source connector logs (%s)\n", r.Name)
			}
			if r.Collection.Destination != "" {
				subTable += fmt.Sprintf("\n# Destination connector logs (%s)\n", r.Name)
			}
		}

		if c.Logs != "" {
			subTable += fmt.Sprintf("\n%s\n", c.Logs)
		}

		if c.Connector.Trace != "" {
			subTable += fmt.Sprintf("\n%s\n", c.Connector.Trace)
		}
	}

	for _, f := range functions {
		if f.Logs != "" {
			subTable += fmt.Sprintf("\n# Function logs (%s)\n\n%s\n", f.Name, f.Logs)
		}
	}

	if deployment != nil && deployment.Status.Details != "" {
		subTable += fmt.Sprintf("\n# Deployment logs (%s)\n\n%s\n", deployment.UUID, deployment.Status.Details)
	}

	return subTable
}
