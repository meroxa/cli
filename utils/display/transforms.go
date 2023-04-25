package display

import (
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

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
				{Text: truncateString(res.Description, 151)}, //nolint:gomnd
			}

			var properties []string
			for _, p := range res.Properties {
				properties = append(properties, p.Name)
			}
			cell := &simpletable.Cell{
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
