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

func AppTable(app *meroxa.Application, resources []*meroxa.Resource, connectors []*meroxa.Connector,
	functions []*meroxa.Function) string {
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

	subTable := extendedResourcesTable(resources, connectors)
	if subTable != "" {
		output += "\n" + subTable
	}
	subTable = extendedFunctionsTable(functions)
	if subTable != "" {
		output += "\n" + subTable
	}
	return output
}

func AppLogsTable(resources []meroxa.ApplicationResource, connectors []*AppExtendedConnector, functions []*meroxa.Function) string {
	subTable := "\tResources:\n"

	var r meroxa.ApplicationResource
	for _, c := range connectors {
		found := false
		for _, resource := range resources {
			if resource.Name.String == c.Connector.ResourceName {
				r = resource
				found = true
				break
			}
		}
		if !found {
			panic("internal error")
		}

		if r.Collection.Source.String != "" {
			subTable += fmt.Sprintf("%s (source: %v)\n", r.Name.String, r.Collection.Source)
		}
		if r.Collection.Destination.String != "" {
			subTable += fmt.Sprintf("%s (destination: %v)\n", r.Name.String, r.Collection.Destination)
		}

		subTable += fmt.Sprintf("%s:\n%s\n", "Trace", c.Connector.Trace)
		subTable += fmt.Sprintf("%s:\n%s\n", "Logs", c.Logs)
	}

	subTable += "\tFunctions:\n"

	for _, f := range functions {
		subTable += fmt.Sprintf("%s Logs:\n%s\n", f.Name, f.Logs)
	}

	return subTable
}
