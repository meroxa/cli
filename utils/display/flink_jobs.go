package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func FlinkJobsTable(flinkJobs []*meroxa.FlinkJob) string {
	if len(flinkJobs) != 0 {
		table := simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
				{Align: simpletable.AlignCenter, Text: "DEPLOYMENT STATE"},
			},
		}

		for _, flinkJob := range flinkJobs {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: flinkJob.UUID},
				{Text: flinkJob.Name},
				{Text: string(flinkJob.Status.State)},
				{Text: string(flinkJob.Status.ManagerDeploymentState)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}

	return ""
}

func PrintFlinkJobsTable(jobs []*meroxa.FlinkJob) {
	fmt.Println(FlinkJobsTable(jobs))
}
