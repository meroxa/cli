package display

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

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
			{Align: simpletable.AlignCenter, Text: string(p.Status.State)},
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
			{Text: string(fun.Status.State)},
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
