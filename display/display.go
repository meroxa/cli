package display

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go"
	"strconv"
	"strings"
)

func JSONPrint(data interface{}) {
	p, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", p)
}

func PrintResourcesTable(resources []*meroxa.Resource) {
	if len(resources) != 0 {
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "ID"},
				{Align: simpletable.AlignCenter, Text: "TYPE"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "URL"},
			},
		}

		for _, res := range resources {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", res.ID)},
				{Text: res.Kind},
				{Text: res.Name},
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
				{Align: simpletable.AlignCenter, Text: "KIND"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "REQUIRED"},
				{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
				{Align: simpletable.AlignCenter, Text: "PROPERTIES"},
			},
		}

		for _, res := range transforms {
			r := []*simpletable.Cell{
				{Text: res.Kind},
				{Text: res.Name},
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
				{Align: simpletable.AlignCenter, Text: "TYPE"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "STREAMS"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
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
				{Text: conn.Kind},
				{Text: conn.Name},
				{Text: streamStr},
				{Text: conn.State},
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
				{Align: simpletable.AlignCenter, Text: "Name"},
			},
		}

		for _, p := range pipelines {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: strconv.Itoa(p.ID)},
				{Align: simpletable.AlignCenter, Text: p.Name},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		fmt.Println(table.String())
	}
}
