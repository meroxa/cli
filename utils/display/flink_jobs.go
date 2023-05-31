package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func FlinkJobsTable(flinkJobs []*meroxa.FlinkJob) string {
	if len(flinkJobs) == 0 {
		return ""
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "UUID"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "ENVIRONMENT"},
			{Align: simpletable.AlignCenter, Text: "STATE"},
			{Align: simpletable.AlignCenter, Text: "DEPLOYMENT STATE"},
		},
	}

	for _, flinkJob := range flinkJobs {
		var env string

		if flinkJob.Environment.Name != "" {
			env = flinkJob.Environment.Name
		} else {
			env = string(meroxa.EnvironmentTypeCommon)
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: flinkJob.UUID},
			{Text: flinkJob.Name},
			{Text: env},
			{Text: string(flinkJob.Status.State)},
			{Text: string(flinkJob.Status.ManagerDeploymentState)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func PrintFlinkJobsTable(jobs []*meroxa.FlinkJob) {
	fmt.Println(FlinkJobsTable(jobs))
}
