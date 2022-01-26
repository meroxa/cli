package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func EndpointsTable(ends []meroxa.Endpoint, hideHeaders bool) string {
	if len(ends) == 0 {
		return ""
	}

	table := simpletable.New()

	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "PROTOCOL"},
				{Align: simpletable.AlignCenter, Text: "STREAM"},
				{Align: simpletable.AlignCenter, Text: "URL"},
				{Align: simpletable.AlignCenter, Text: "READY"},
			},
		}
	}

	for _, end := range ends {
		var u string
		switch end.Protocol {
		case meroxa.EndpointProtocolHttp:
			host, err := url.ParseRequestURI(end.Host)
			if err != nil {
				continue
			}
			host.User = url.UserPassword(end.BasicAuthUsername, end.BasicAuthPassword)
			u = host.String()
		case meroxa.EndpointProtocolGrpc:
			u = fmt.Sprintf("host=%s username=%s password=%s", end.Host, end.BasicAuthUsername, end.BasicAuthPassword)
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: end.Name},
			{Text: string(end.Protocol)},
			{Text: end.Stream},
			{Text: u},
			{Text: strings.Title(strconv.FormatBool(end.Ready))},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)

	return table.String()
}

func ResourceTable(res *meroxa.Resource) string {
	tunnel := "N/A"
	if res.SSHTunnel != nil {
		tunnel = "SSH"
	}

	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "ID:"},
			{Text: fmt.Sprintf("%d", res.ID)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: res.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(res.Type)},
		},
		{
			{Align: simpletable.AlignRight, Text: "URL:"},
			{Text: res.URL},
		},
		{
			{Align: simpletable.AlignRight, Text: "Tunnel:"},
			{Text: tunnel},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: strings.Title(string(res.Status.State))},
		},
	}

	if d := res.Status.Details; d != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State details:"},
			{Text: strings.Title(d)},
		})
	}

	if res.Environment != nil {
		if e := res.Environment.UUID; e != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: e},
			})
		}

		if e := res.Environment.Name; e != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: e},
			})
		}
	} else {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Environment Name:"},
			{Text: string(meroxa.EnvironmentTypeCommon)},
		})
	}

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func PipelineTable(p *meroxa.Pipeline) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: p.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "ID:"},
			{Text: fmt.Sprintf("%d", p.ID)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: p.Name},
		},
	}

	if p.Environment != nil {
		if pU := p.Environment.UUID; pU != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: pU},
			})
		}
		if pN := p.Environment.Name; pN != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: pN},
			})
		}
	} else {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Environment Name:"},
			{Text: string(meroxa.EnvironmentTypeCommon)},
		})
	}

	mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "State:"},
		{Text: strings.Title(string(p.State))},
	})

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func PrintPipelineTable(pipeline *meroxa.Pipeline) {
	fmt.Println(PipelineTable(pipeline))
}

func ResourcesTable(resources []*meroxa.Resource, hideHeaders bool) string {
	if len(resources) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
					{Align: simpletable.AlignCenter, Text: "URL"},
					{Align: simpletable.AlignCenter, Text: "TUNNEL"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, res := range resources {
			tunnel := "N/A"
			if res.SSHTunnel != nil {
				tunnel = "SSH"
			}

			var env string

			if res.Environment != nil && res.Environment.Name != "" {
				env = res.Environment.Name
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", res.ID)},
				{Text: res.Name},
				{Text: string(res.Type)},
				{Text: env},
				{Text: res.URL},
				{Align: simpletable.AlignCenter, Text: tunnel},
				{Align: simpletable.AlignCenter, Text: strings.Title(string(res.Status.State))},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
}

func PrintResourcesTable(resources []*meroxa.Resource, hideHeaders bool) {
	fmt.Println(ResourcesTable(resources, hideHeaders))
}

func TransformsTable(transforms []*meroxa.Transform, hideHeaders bool) string {
	if len(transforms) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "REQUIRED"},
					{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
					{Align: simpletable.AlignCenter, Text: "PROPERTIES"},
				},
			}
		}

		for _, res := range transforms {
			r := []*simpletable.Cell{
				{Text: res.Name},
				{Text: res.Type},
				{Text: strconv.FormatBool(res.Required)},
				{Text: truncateString(res.Description, 151)}, // nolint:gomnd
			}

			var properties []string
			for _, p := range res.Properties {
				properties = append(properties, p.Name)
			}
			var cell = &simpletable.Cell{
				Text: strings.Join(properties, ","),
			}

			r = append(r, cell)
			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
}

func ConnectorTable(connector *meroxa.Connector) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: connector.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "ID:"},
			{Text: fmt.Sprintf("%d", connector.ID)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: connector.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(connector.Type)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Streams:"},
			{Text: formatStreams(connector.Streams)},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(connector.State)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Pipeline:"},
			{Text: connector.PipelineName},
		},
	}

	if connector.Trace != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Trace:"},
			{Text: connector.Trace},
		})
	}
	if connector.Environment != nil {
		if envUUID := connector.Environment.UUID; envUUID != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: envUUID},
			})
		}
		if envName := connector.Environment.Name; envName != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: envName},
			})
		}
	} else {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Environment Name:"},
			{Text: string(meroxa.EnvironmentTypeCommon)},
		})
	}

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func ConnectorsTable(connectors []*meroxa.Connector, hideHeaders bool) string {
	if len(connectors) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "STREAMS"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
					{Align: simpletable.AlignCenter, Text: "PIPELINE"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
				},
			}
		}

		for _, conn := range connectors {
			var env string

			if conn.Environment != nil && conn.Environment.Name != "" {
				env = conn.Environment.Name
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			streamStr := formatStreams(conn.Streams)
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: conn.UUID},
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", conn.ID)},
				{Text: conn.Name},
				{Text: string(conn.Type)},
				{Text: streamStr},
				{Text: string(conn.State)},
				{Text: conn.PipelineName},
				{Text: env},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
}

func formatStreams(ss map[string]interface{}) string {
	var streamStr string

	if streamInput, ok := ss["input"]; ok {
		streamStr += "(input) "

		streamInterface := streamInput.([]interface{})
		for i, v := range streamInterface {
			streamStr += fmt.Sprintf("%v", v)

			if i < len(streamInterface)-1 {
				streamStr += ", "
			}
		}
	}

	if streamOutput, ok := ss["output"]; ok {
		streamStr += "(output) "

		streamInterface := streamOutput.([]interface{})
		for i, v := range streamInterface {
			streamStr += fmt.Sprintf("%v", v)

			if i < len(streamInterface)-1 {
				streamStr += ", "
			}
		}
	}
	return streamStr
}

func PrintConnectorsTable(connectors []*meroxa.Connector, hideHeaders bool) {
	fmt.Println(ConnectorsTable(connectors, hideHeaders))
}

func ResourceTypesTable(types []string, hideHeaders bool) string {
	table := simpletable.New()

	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "TYPES"},
			},
		}
	}

	for _, t := range types {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: t},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func PrintResourceTypesTable(types []string, hideHeaders bool) {
	fmt.Println(ResourceTypesTable(types, hideHeaders))
}

func PipelinesTable(pipelines []*meroxa.Pipeline, hideHeaders bool) string {
	if len(pipelines) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, p := range pipelines {
			var env string

			if p.Environment != nil && p.Environment.Name != "" {
				env = p.Environment.Name
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: p.UUID},
				{Align: simpletable.AlignRight, Text: strconv.Itoa(p.ID)},
				{Align: simpletable.AlignCenter, Text: p.Name},
				{Align: simpletable.AlignCenter, Text: env},
				{Align: simpletable.AlignCenter, Text: string(p.State)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}
	return ""
}

func PrintPipelinesTable(pipelines []*meroxa.Pipeline, hideHeaders bool) {
	fmt.Println(PipelinesTable(pipelines, hideHeaders))
}

func FunctionsTable(funs []*meroxa.Function, hideHeaders bool) string {
	if len(funs) == 0 {
		return ""
	}

	table := simpletable.New()
	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "INPUT STREAM"},
				{Align: simpletable.AlignCenter, Text: "OUTPUT STREAM"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
				{Align: simpletable.AlignCenter, Text: "PIPELINE"},
			},
		}
	}

	for _, p := range funs {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: p.UUID},
			{Align: simpletable.AlignCenter, Text: p.Name},
			{Align: simpletable.AlignCenter, Text: p.InputStream},
			{Align: simpletable.AlignCenter, Text: p.OutputStream},
			{Align: simpletable.AlignCenter, Text: p.Status.State},
			{Align: simpletable.AlignCenter, Text: p.Pipeline.Name},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func FunctionTable(fun *meroxa.Function) string {
	envVars := []string{}
	for k, v := range fun.EnvVars {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: fun.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: fun.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Input Stream:"},
			{Text: fun.InputStream},
		},
		{
			{Align: simpletable.AlignRight, Text: "Output Stream:"},
			{Text: fun.OutputStream},
		},
		{
			{Align: simpletable.AlignRight, Text: "Image:"},
			{Text: fun.Image},
		},
		{
			{Align: simpletable.AlignRight, Text: "Command:"},
			{Text: strings.Join(fun.Command, " ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Arguments:"},
			{Text: strings.Join(fun.Args, " ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Environment Variables:"},
			{Text: strings.Join(envVars, "\n")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Pipeline:"},
			{Text: fun.Pipeline.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: strings.Title(fun.Status.State)},
		},
	}
	mainTable.SetStyle(simpletable.StyleCompact)
	table := mainTable.String()

	details := fun.Status.Details
	if details == "" {
		return table
	}

	return fmt.Sprintf("%s\n\n%s", table, details)
}

func EnvironmentsTable(environments []*meroxa.Environment, hideHeaders bool) string {
	if len(environments) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "PROVIDER"},
					{Align: simpletable.AlignCenter, Text: "REGION"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, p := range environments {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: p.UUID},
				{Align: simpletable.AlignCenter, Text: p.Name},
				{Align: simpletable.AlignCenter, Text: string(p.Type)},
				{Align: simpletable.AlignCenter, Text: string(p.Provider)},
				{Align: simpletable.AlignCenter, Text: string(p.Region)},
				{Align: simpletable.AlignCenter, Text: string(p.Status.State)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}
	return ""
}

func EnvironmentTable(environment *meroxa.Environment) string {
	mainTable := simpletable.New()

	bytes, _ := json.Marshal(&environment.Status)

	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: environment.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: environment.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Provider:"},
			{Text: string(environment.Provider)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Region:"},
			{Text: string(environment.Region)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(environment.Type)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Status:"},
			{Text: string(bytes)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: environment.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: environment.UpdatedAt.String()},
		},
	}

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func PrintEnvironmentsTable(environments []*meroxa.Environment, hideHeaders bool) {
	fmt.Println(EnvironmentsTable(environments, hideHeaders))
}

func truncateString(oldString string, l int) string {
	str := oldString

	if len(oldString) > l {
		str = oldString[:l] + "..."
	}

	return str
}
