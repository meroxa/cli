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
				// {Align: simpletable.AlignCenter, Text: "INPUTS TREAMS"},
				// {Align: simpletable.AlignCenter, Text: "OUTPUT STREAMS"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
				// {Align: simpletable.AlignCenter, Text: "LIFECYCLE STATE"},
				// {Align: simpletable.AlignCenter, Text: "RECONCILIATION STATE"},
				{Align: simpletable.AlignCenter, Text: "DEPLOYMENT STATE"},
				// {Align: simpletable.AlignCenter, Text: "CREATED"},
				// {Align: simpletable.AlignCenter, Text: "UPDATED"},
			},
		}

		for _, flinkJob := range flinkJobs {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: flinkJob.UUID},
				{Text: flinkJob.Name},
				// {Text: fmt.Sprintf("%v", flinkJob.InputStreams)},
				// {Text: fmt.Sprintf("%v", flinkJob.OutputStreams)},
				{Text: string(flinkJob.Status.State)},
				// {Text: string(flinkJob.Status.LifecycleState)},
				// {Text: string(flinkJob.Status.ReconciliationState)},
				{Text: string(flinkJob.Status.ManagerDeploymentState)},
				// {Text: flinkJob.CreatedAt.String()},
				// {Text: flinkJob.UpdatedAt.String()},
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

/* func FlinkJobTable(flinkJob *meroxa.FlinkJob) string {
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
			{Align: simpletable.AlignRight, Text: "InputStreams:"},
			{Text: fmt.Sprintf("%v", flinkJob.InputStreams)},
		},
		{
			{Align: simpletable.AlignRight, Text: "OutputStreams:"},
			{Text: fmt.Sprintf("%v", flinkJob.OutputStreams)},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(flinkJob.Status.State)},
		},
		{
			{Align: simpletable.AlignRight, Text: "LifeCycleState:"},
			{Text: string(flinkJob.Status.LifecycleState)},
		},
		{
			{Align: simpletable.AlignRight, Text: "ReconciliationState:"},
			{Text: string(flinkJob.Status.ReconciliationState)},
		},
		{
			{Align: simpletable.AlignRight, Text: "ManagerDeploymentState:"},
			{Text: string(flinkJob.Status.ManagerDeploymentState)},
		},
		{
			{Align: simpletable.AlignRight, Text: "CreatedAt:"},
			{Text: flinkJob.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "UpdatedAt:"},
			{Text: flinkJob.UpdatedAt.String()},
		},
	}

	if flinkJob.Status.Details != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State details:"},
			{Text: flinkJob.Status.Details},
		})
	}

	if flinkJob.Environment != nil {
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
} */
