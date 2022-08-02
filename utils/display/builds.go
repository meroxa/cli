package display

import (
	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func BuildTable(build *meroxa.Build) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: build.Uuid},
		},
		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: build.CreatedAt},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: build.UpdatedAt},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: build.Status.State},
		},
	}
	if build.Status.Details != "" {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Status Details:"},
			{Text: build.Status.Details},
		}
		mainTable.Body.Cells = append(mainTable.Body.Cells, r)
	}
	mainTable.SetStyle(simpletable.StyleCompact)
	return mainTable.String()
}
