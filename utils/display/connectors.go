package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func formatStreams(ss map[string]interface{}) string {
	var streamStr string

	if streamInput, ok := ss["input"]; ok {
		streamStr += "(input) "

		streamInterface := streamInput.([]interface{})
		for i, v := range streamInterface {
			streamStr += fmt.Sprintf("%v", v)

			if i < len(streamInterface)-1 {
				streamStr += ", "
			}
		}
	}

	if streamOutput, ok := ss["output"]; ok {
		streamStr += "(output) "

		streamInterface := streamOutput.([]interface{})
		for i, v := range streamInterface {
			streamStr += fmt.Sprintf("%v", v)

			if i < len(streamInterface)-1 {
				streamStr += ", "
			}
		}
	}
	return streamStr
}

func ConnectorTable(connector *meroxa.Connector) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: connector.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: connector.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(connector.Type)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Streams:"},
			{Text: formatStreams(connector.Streams)},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(connector.State)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Pipeline:"},
			{Text: connector.PipelineName},
		},
	}

	if connector.Trace != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Trace:"},
			{Text: connector.Trace},
		})
	}
	if connector.Environment != nil {
		if connector.Environment.UUID.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: connector.Environment.UUID.String},
			})
		}
		if connector.Environment.Name.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: connector.Environment.Name.String},
			})
		}
	} else {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Environment Name:"},
			{Text: string(meroxa.EnvironmentTypeCommon)},
		})
	}

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func ConnectorsTable(connectors []*meroxa.Connector, hideHeaders bool) string {
	if len(connectors) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "STREAMS"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
					{Align: simpletable.AlignCenter, Text: "PIPELINE"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
				},
			}
		}

		for _, conn := range connectors {
			var env string

			if conn.Environment != nil && conn.Environment.Name.Valid {
				env = conn.Environment.Name.String
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			streamStr := formatStreams(conn.Streams)
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: conn.UUID},
				{Text: conn.Name},
				{Text: string(conn.Type)},
				{Text: streamStr},
				{Text: string(conn.State)},
				{Text: conn.PipelineName},
				{Text: env},
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
