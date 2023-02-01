package display

import (
	"fmt"
	"sort"

	"github.com/alexeyco/simpletable"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func ResourceTable(res *meroxa.Resource) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: res.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: res.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(res.Type)},
		},
	}
	if res.URL != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "URL:"},
			{Text: res.URL},
		})
	}
	if res.SSHTunnel != nil {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Tunnel:"},
			{Text: "SSH"},
		})
	}

	mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: "State:"},
		{Text: string(res.Status.State)},
	})
	if res.Status.Details != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State details:"},
			{Text: res.Status.Details},
		})
	}

	if res.Environment != nil {
		if res.Environment.UUID != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: res.Environment.UUID},
			})
		}

		if res.Environment.Name != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: res.Environment.Name},
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

func ResourcesTable(resources []*meroxa.Resource, hideHeaders bool) string {
	if len(resources) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
					{Align: simpletable.AlignCenter, Text: "URL"},
					{Align: simpletable.AlignCenter, Text: "TUNNEL"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, res := range resources {
			tunnel := "N/A"
			if res.SSHTunnel != nil {
				tunnel = "SSH"
			}

			var env string

			if res.Environment != nil && res.Environment.Name != "" {
				env = res.Environment.Name
			} else {
				env = string(meroxa.EnvironmentTypeCommon)
			}

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: res.UUID},
				{Text: res.Name},
				{Text: string(res.Type)},
				{Text: env},
				{Text: res.URL},
				{Align: simpletable.AlignCenter, Text: tunnel},
				{Align: simpletable.AlignCenter, Text: string(res.Status.State)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
}

func PrintResourcesTable(resources []*meroxa.Resource, hideHeaders bool) {
	fmt.Println(ResourcesTable(resources, hideHeaders))
}

func ResourceTypesTable(types []meroxa.ResourceType, hideHeaders bool) string {
	gaResourceNames := []string{}
	gaResourceTypes := make(map[string]string)
	betaResourceNames := []string{}
	betaResourceTypes := make(map[string]string)

	for _, t := range types {
		val := fmt.Sprintf("%s", t.FormConfig[meroxa.ResourceTypeFormConfigHumanReadableKey])
		if val == "" {
			continue
		}
		if t.ReleaseStage == meroxa.ResourceTypeReleaseStageGA {
			gaResourceNames = append(gaResourceNames, val)
			gaResourceTypes[val] = t.Name
		} else if t.ReleaseStage == meroxa.ResourceTypeReleaseStageBeta {
			betaResourceNames = append(betaResourceNames, val)
			betaResourceTypes[val] = t.Name
		}
	}
	sort.Strings(gaResourceNames)
	sort.Strings(betaResourceNames)

	table := simpletable.New()

	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "TYPE"},
				{Align: simpletable.AlignCenter, Text: "RELEASE STAGE"},
			},
		}
	}

	for _, t := range gaResourceNames {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: t},
			{Align: simpletable.AlignLeft, Text: gaResourceTypes[t]},
			{Align: simpletable.AlignLeft, Text: string(meroxa.ResourceTypeReleaseStageGA)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	for _, t := range betaResourceNames {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: t},
			{Align: simpletable.AlignLeft, Text: betaResourceTypes[t]},
			{Align: simpletable.AlignLeft, Text: string(meroxa.ResourceTypeReleaseStageBeta)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func PrintResourceTypesTable(types []meroxa.ResourceType, hideHeaders bool) {
	fmt.Println(ResourceTypesTable(types, hideHeaders))
}
