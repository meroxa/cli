package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go"
)

func JSONPrint(data interface{}) {
	p, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", p)
}

func EndpointsTable(ends []meroxa.Endpoint) string {
	if len(ends) == 0 {
		return ""
	}

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "PROTOCOL"},
			{Align: simpletable.AlignCenter, Text: "STREAM"},
			{Align: simpletable.AlignCenter, Text: "URL"},
			{Align: simpletable.AlignCenter, Text: "READY"},
		},
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
func PrintEndpointsTable(ends []meroxa.Endpoint) {
	fmt.Println(EndpointsTable(ends))
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

func ResourcesTable(resources []*meroxa.Resource) string {
	if len(resources) != 0 {
		table := simpletable.New()
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

func PrintResourcesTable(resources []*meroxa.Resource) {
	fmt.Println(ResourcesTable(resources))
}

func PrintTransformsTable(transforms []*meroxa.Transform) {
	TransformsTable(transforms)
}

func TransformsTable(transforms []*meroxa.Transform) string {
	if len(transforms) != 0 {
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "TYPE"},
				{Align: simpletable.AlignCenter, Text: "REQUIRED"},
				{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
				{Align: simpletable.AlignCenter, Text: "PROPERTIES"},
			},
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

func ConnectorsTable(connectors []*meroxa.Connector) string {
	if len(connectors) != 0 {
		table := simpletable.New()
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

func PrintConnectorsTable(connectors []*meroxa.Connector) {
	fmt.Println(ConnectorsTable(connectors))
}

func ResourceTypesTable(types []string) string {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "TYPES"},
		},
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

func PrintResourceTypesTable(types []string) {
	fmt.Println(ResourceTypesTable(types))
}

func PipelinesTable(pipelines []*meroxa.Pipeline) string {
	if len(pipelines) != 0 {
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "ID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
			},
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

func PrintPipelinesTable(pipelines []*meroxa.Pipeline) {
	fmt.Println(PipelinesTable(pipelines))
}
