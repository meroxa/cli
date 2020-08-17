package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go"
	"strconv"
	"strings"
)

func prettyPrint(section string, data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("== %s ==\n", strings.ToTitle(section))
	fmt.Printf("%s \n", p)
}

func printResourcesTable(resources []*meroxa.Resource) {
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

func appendCell(cells []*simpletable.Cell, text string) []*simpletable.Cell {
	cells = append(cells, &simpletable.Cell{
		Align: simpletable.AlignCenter,
		Text:  text,
	})
	return cells
}

func printConnectionsTable(connections []*meroxa.Connector) {
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

	for _, conn := range connections {
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

func printResourceTypesTable(types []string) {
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

func printPipelinesTable(pipelines []*meroxa.Pipeline) {
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
