package display

import (
	"fmt"
	"net/url"
	"strconv"

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
			{Text: strconv.FormatBool(end.Ready)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)

	return table.String()
}
