package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func PipelineTable(p *meroxa.Pipeline) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: p.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: p.Name},
		},
	}

	if p.Environment != nil {
		if p.Environment.UUID.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: p.Environment.UUID.String},
			})
		}
		if p.Environment.Name.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: p.Environment.Name.String},
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
		{Text: string(p.State)},
	})

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func PrintPipelineTable(pipeline *meroxa.Pipeline) {
	fmt.Println(PipelineTable(pipeline))
}

func PipelinesTable(pipelines []*meroxa.Pipeline, hideHeaders bool) string {
	if len(pipelines) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, p := range pipelines {
			var env string

			if p.Environment != nil && p.Environment.Name.Valid {
				env = p.Environment.Name.String
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: p.UUID},
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
