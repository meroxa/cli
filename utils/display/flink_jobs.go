package display

import (
	"fmt"
	"strings"

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
			{Align: simpletable.AlignCenter, Text: "LIFECYCLE STATE"},
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
			{Text: string(flinkJob.Status.LifecycleState)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func PrintFlinkJobsTable(jobs []*meroxa.FlinkJob) {
	fmt.Println(FlinkJobsTable(jobs))
}

func FlinkJobTable(flinkJob *meroxa.FlinkJob) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: flinkJob.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: flinkJob.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: flinkJob.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: flinkJob.UpdatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Input Streams:"},
			{Text: strings.Join(flinkJob.InputStreams, ", ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Output Streams:"},
			{Text: strings.Join(flinkJob.OutputStreams, ", ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Lifecycle State:"},
			{Text: string(flinkJob.Status.LifecycleState)},
		},
	}
	if flinkJob.Status.State != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Job State:"},
			{Text: flinkJob.Status.State},
		})
	}

	if flinkJob.Status.ReconciliationState != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Reconciliation State:"},
			{Text: string(flinkJob.Status.ReconciliationState)},
		})
	}

	if flinkJob.Status.ManagerDeploymentState != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Manager Deployment State:"},
			{Text: string(flinkJob.Status.ManagerDeploymentState)},
		})
	}

	if flinkJob.Status.Details != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State Details:"},
			{Text: flinkJob.Status.Details},
		})
	}

	if flinkJob.Environment.Name != "" {
		if flinkJob.Environment.UUID != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: flinkJob.Environment.UUID},
			})
		}

		if flinkJob.Environment.Name != "" {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: flinkJob.Environment.Name},
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
