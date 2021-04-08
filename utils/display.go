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

func PrintEndpointsTable(ends []meroxa.Endpoint) {
	if len(ends) == 0 {
		return
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
	fmt.Println(table.String())
}

func PrintResourcesTable(resources []*meroxa.Resource) {
	if len(resources) != 0 {
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "ID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "TYPE"},
				{Align: simpletable.AlignCenter, Text: "URL"},
			},
		}

		for _, res := range resources {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", res.ID)},
				{Text: res.Name},
				{Text: res.Type},
				{Text: res.URL},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		fmt.Println(table.String())
	}
}

func PrintTransformsTable(transforms []*meroxa.Transform) {
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
		fmt.Println(table.String())
	}
}

func PrintConnectorsTable(connectors []*meroxa.Connector) {
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
				{Text: fmt.Sprintf("%s", conn.PipelineName)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		fmt.Println(table.String())
	}
}

func PrintResourceTypesTable(types []string) {
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
	fmt.Println(table.String())
}

func PrintPipelinesTable(pipelines []*meroxa.Pipeline) {
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
		fmt.Println(table.String())
	}
}
