package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go"
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
		case "HTTP":
			host, err := url.ParseRequestURI(end.Host)
			if err != nil {
				continue
			}
			host.User = url.UserPassword(end.BasicAuthUsername, end.BasicAuthPassword)
			u = host.String()
		case "GRPC":
			u = fmt.Sprintf("host=%s username=%s password=%s", end.Host, end.BasicAuthUsername, end.BasicAuthPassword)
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: end.Name},
			{Text: end.Protocol},
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
			{Text: res.Type},
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
			{Text: strings.Title(res.Status.State)},
		},
	}

	if d := res.Status.Details; d != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State details:"},
			{Text: strings.Title(d)},
		})
	}
	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
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

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", res.ID)},
				{Text: res.Name},
				{Text: res.Type},
				{Text: res.URL},
				{Align: simpletable.AlignCenter, Text: tunnel},
				{Align: simpletable.AlignCenter, Text: strings.Title(res.Status.State)},
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
				{Text: strings.ReplaceAll(res.Description, ". ", ". \n")},
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

func ConnectorsTable(connectors []*meroxa.Connector, hideHeaders bool) string {
	if len(connectors) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "STREAMS"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
					{Align: simpletable.AlignCenter, Text: "PIPELINE"},
				},
			}
		}

		for _, conn := range connectors {
			var streamStr string

			if streamInput, ok := conn.Streams["input"]; ok {
				streamStr += "input:\n"

				streamInterface := streamInput.([]interface{})
				for _, v := range streamInterface {
					streamStr += fmt.Sprintf("%v\n", v)
				}
			}

			if streamOutput, ok := conn.Streams["output"]; ok {
				streamStr += "output:\n"

				streamInterface := streamOutput.([]interface{})
				for _, v := range streamInterface {
					streamStr += fmt.Sprintf("%v\n", v)
				}
			}
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", conn.ID)},
				{Text: conn.Name},
				{Text: conn.Type},
				{Text: streamStr},
				{Text: conn.State},
				{Text: conn.PipelineName},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
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
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, p := range pipelines {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: strconv.Itoa(p.ID)},
				{Align: simpletable.AlignCenter, Text: p.Name},
				{Align: simpletable.AlignCenter, Text: p.State},
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
