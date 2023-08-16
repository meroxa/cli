package display

import (
	"fmt"
	"strings"
	"time"

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

	for _, app := range apps {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: app.UUID},
			{Align: simpletable.AlignLeft, Text: app.Name},
			{Align: simpletable.AlignLeft, Text: app.Language},
			{Align: simpletable.AlignLeft, Text: app.GitSha},
			{Align: simpletable.AlignLeft, Text: string(app.Status.State)},
		}

		if app.Environment != nil && app.Environment.Name != "" {
			r = append(r, &simpletable.Cell{Align: simpletable.AlignLeft, Text: app.Environment.Name})
		} else {
			r = append(r, &simpletable.Cell{Align: simpletable.AlignLeft, Text: string(meroxa.EnvironmentTypeCommon)})
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	if !hideHeaders {
		cells := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "UUID"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "LANGUAGE"},
			{Align: simpletable.AlignCenter, Text: "GIT SHA"},
			{Align: simpletable.AlignCenter, Text: "STATE"},
			{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
		}

		table.Header = &simpletable.Header{
			Cells: cells,
		}
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

	if app.Environment != nil {
		if _, ok := app.Environment.GetNameOrUUID(); ok == nil {
			subTable = appEnvironmentTable(app.Environment)
			if subTable != "" {
				output += "\n" + subTable
			}
		}
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

	for _, r := range resources {
		var status string
		t := "source"
		if r.Collection.Destination == "true" {
			t = "destination"
		}
		for _, c := range connectors {
			if r.UUID == c.ResourceUUID && isMatching(r.Collection, c.ResourceType) {
				status = c.Status
				break
			}
		}

		subTable += fmt.Sprintf("\t    %s (%s)\n", r.Name, t)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", r.UUID)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "Type", r.ResourceType)
		if status != "" {
			subTable += fmt.Sprintf("\t\t%5s:   %s\n", "State", status)
		}
		status = ""
	}

	return subTable
}

func appEnvironmentTable(env *meroxa.EntityIdentifier) string {
	subTable := "\tEnvironment\n"

	subTable += fmt.Sprintf("\t    %s\n", env.Name)
	subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", env.UUID)

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

func AppLogsTable(appLogs *meroxa.ApplicationLogs) string {
	var subTable string

	for key, logs := range appLogs.ConnectorLogs {
		subTable += fmt.Sprintf("\n# Logs for %s resource\n\n%s\n", key, logs)
	}

	for key, logs := range appLogs.FunctionLogs {
		subTable += fmt.Sprintf("\n# Logs for %s function\n\n%s\n", key, logs)
	}

	for key, logs := range appLogs.DeploymentLogs {
		subTable += fmt.Sprintf("\n# Logs for %s deployment\n\n%s\n", key, logs)
	}

	return subTable
}

func AppLogsTableV2(appLogs *meroxa.Logs) string {
	var subTable string

	for _, l := range appLogs.Data {
		subTable += fmt.Sprintf("[%s]\t%s\t%q\n", l.Timestamp.Format(time.RFC3339), l.Source, l.Log)
	}

	return subTable
}
