package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

func ResourceTable(res *meroxa.Resource) string {
	tunnel := "N/A"
	if res.SSHTunnel != nil {
		tunnel = "SSH"
	}

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
		{
			{Align: simpletable.AlignRight, Text: "URL:"},
			{Text: res.URL},
		},
		{
			{Align: simpletable.AlignRight, Text: "Tunnel:"},
			{Text: tunnel},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(res.Status.State)},
		},
	}

	if res.Status.Details != "" {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "State details:"},
			{Text: res.Status.Details},
		})
	}

	if res.Environment != nil {
		if res.Environment.UUID.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment UUID:"},
				{Text: res.Environment.UUID.String},
			})
		}

		if res.Environment.Name.Valid {
			mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: "Environment Name:"},
				{Text: res.Environment.Name.String},
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

			if res.Environment != nil && res.Environment.Name.Valid {
				env = res.Environment.Name.String
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
				{Text: truncateString(res.Description, 151)}, // nolint:gomnd
			}

			var properties []string
			for _, p := range res.Properties {
				properties = append(properties, p.Name)
			}
			var cell = &simpletable.Cell{
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

func PrintConnectorsTable(connectors []*meroxa.Connector, hideHeaders bool) {
	fmt.Println(ConnectorsTable(connectors, hideHeaders))
}

func ResourceTypesTable(types []string, hideHeaders bool) string {
	table := simpletable.New()

	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "TYPES"},
			},
		}
	}

	for _, t := range types {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: t},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func PrintResourceTypesTable(types []string, hideHeaders bool) {
	fmt.Println(ResourceTypesTable(types, hideHeaders))
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

func FunctionsTable(funs []*meroxa.Function, hideHeaders bool) string {
	if len(funs) == 0 {
		return ""
	}

	table := simpletable.New()
	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "INPUT STREAM"},
				{Align: simpletable.AlignCenter, Text: "OUTPUT STREAM"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
				{Align: simpletable.AlignCenter, Text: "PIPELINE"},
			},
		}
	}

	for _, p := range funs {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: p.UUID},
			{Align: simpletable.AlignCenter, Text: p.Name},
			{Align: simpletable.AlignCenter, Text: p.InputStream},
			{Align: simpletable.AlignCenter, Text: p.OutputStream},
			{Align: simpletable.AlignCenter, Text: string(p.Status.State)},
			{Align: simpletable.AlignCenter, Text: p.Pipeline.Name},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func FunctionTable(fun *meroxa.Function) string {
	envVars := []string{}
	for k, v := range fun.EnvVars {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: fun.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: fun.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Input Stream:"},
			{Text: fun.InputStream},
		},
		{
			{Align: simpletable.AlignRight, Text: "Output Stream:"},
			{Text: fun.OutputStream},
		},
		{
			{Align: simpletable.AlignRight, Text: "Image:"},
			{Text: fun.Image},
		},
		{
			{Align: simpletable.AlignRight, Text: "Command:"},
			{Text: strings.Join(fun.Command, " ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Arguments:"},
			{Text: strings.Join(fun.Args, " ")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Environment Variables:"},
			{Text: strings.Join(envVars, "\n")},
		},
		{
			{Align: simpletable.AlignRight, Text: "Pipeline:"},
			{Text: fun.Pipeline.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(fun.Status.State)},
		},
	}
	mainTable.SetStyle(simpletable.StyleCompact)
	table := mainTable.String()

	details := fun.Status.Details
	if details == "" {
		return table
	}

	return fmt.Sprintf("%s\n\n%s", table, details)
}

func EnvironmentsTable(environments []*meroxa.Environment, hideHeaders bool) string {
	if len(environments) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
					{Align: simpletable.AlignCenter, Text: "TYPE"},
					{Align: simpletable.AlignCenter, Text: "PROVIDER"},
					{Align: simpletable.AlignCenter, Text: "REGION"},
					{Align: simpletable.AlignCenter, Text: "STATE"},
				},
			}
		}

		for _, p := range environments {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: p.UUID},
				{Align: simpletable.AlignCenter, Text: p.Name},
				{Align: simpletable.AlignCenter, Text: string(p.Type)},
				{Align: simpletable.AlignCenter, Text: string(p.Provider)},
				{Align: simpletable.AlignCenter, Text: string(p.Region)},
				{Align: simpletable.AlignCenter, Text: string(p.Status.State)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}
	return ""
}

// nolint:funlen
func EnvironmentTable(environment *meroxa.Environment) string {
	mainTable := simpletable.New()

	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: environment.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: environment.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Provider:"},
			{Text: string(environment.Provider)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Region:"},
			{Text: string(environment.Region)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Type:"},
			{Text: string(environment.Type)},
		},

		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: environment.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: environment.UpdatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Environment Status:"},
			{Text: string(environment.Status.State)},
		},
	}

	if environment.Status.Details != "" {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Environment Status Details:"},
			{Text: string(environment.Status.State)},
		}
		mainTable.Body.Cells = append(mainTable.Body.Cells, r)
	}

	mainTable.SetStyle(simpletable.StyleCompact)
	str := mainTable.String()

	if environment.Status.PreflightDetails != nil {
		preflightTable := simpletable.New()
		preflightTable.Body.Cells = [][]*simpletable.Cell{
			{
				{Align: simpletable.AlignRight, Text: "				Preflight Checks:"},
				{Text: ""},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EC2 Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.EC2, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EKS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.EKS, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS IAM Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.IAM, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS KMS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.KMS, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS MKS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.MSK, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS S3 Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.S3, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS ServiceQuotas Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.ServiceQuotas, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS CloudFormation Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.Cloudformation, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS Cloudwatch Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.Cloudwatch, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EIP Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.EIP},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS NAT Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.NAT},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS VPC Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.VPC},
			},
		}
		preflightTable.SetStyle(simpletable.StyleCompact)
		str += "\n" + preflightTable.String()
	}

	return str
}

func EnvironmentPreflightTable(environment *meroxa.Environment) string {
	if environment.Status.PreflightDetails != nil {
		preflightTable := simpletable.New()
		preflightTable.Body.Cells = [][]*simpletable.Cell{
			{
				{Align: simpletable.AlignRight, Text: "				Preflight Checks:"},
				{Text: ""},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EC2 Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.EC2, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EKS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.EKS, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS IAM Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.IAM, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS KMS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.KMS, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS MKS Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.MSK, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS S3 Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.S3, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS ServiceQuotas Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.ServiceQuotas, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS CloudFormation Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.Cloudformation, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS Cloudwatch Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.Cloudwatch, " ; ")},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS EIP Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.EIP},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS NAT Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.NAT},
			},
			{
				{Align: simpletable.AlignRight, Text: "AWS VPC Limits Status:"},
				{Text: environment.Status.PreflightDetails.PreflightLimits.VPC},
			},
		}
		preflightTable.SetStyle(simpletable.StyleCompact)
		return preflightTable.String()
	}
	return ""
}

func PrintEnvironmentsTable(environments []*meroxa.Environment, hideHeaders bool) {
	fmt.Println(EnvironmentsTable(environments, hideHeaders))
}

func AppsTable(apps []*meroxa.Application, hideHeaders bool) string {
	if len(apps) == 0 {
		return ""
	}

	table := simpletable.New()
	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "LANGUAGE"},
				{Align: simpletable.AlignCenter, Text: "GIT SHA"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
			},
		}
	}

	for _, app := range apps {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: app.UUID},
			{Align: simpletable.AlignCenter, Text: app.Name},
			{Align: simpletable.AlignCenter, Text: app.Language},
			{Align: simpletable.AlignCenter, Text: app.GitSha},
			{Align: simpletable.AlignCenter, Text: string(app.Status.State)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

func AppTable(app *meroxa.Application, resources []*meroxa.Resource, connectors []*meroxa.Connector,
	functions []*meroxa.Function) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: app.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: app.Name},
		},
		{
			{Align: simpletable.AlignRight, Text: "Language:"},
			{Text: app.Language},
		},
		{
			{Align: simpletable.AlignRight, Text: "Git SHA:"},
			{Text: strings.TrimSpace(app.GitSha)},
		},
		{
			{Align: simpletable.AlignRight, Text: "Created At:"},
			{Text: app.CreatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "Updated At:"},
			{Text: app.UpdatedAt.String()},
		},
		{
			{Align: simpletable.AlignRight, Text: "State:"},
			{Text: string(app.Status.State)},
		},
	}
	mainTable.SetStyle(simpletable.StyleCompact)
	output := mainTable.String()

	subTable := extendedResourcesTable(resources, connectors)
	if subTable != "" {
		output += "\n" + subTable
	}
	subTable = extendedFunctionsTable(functions)
	if subTable != "" {
		output += "\n" + subTable
	}
	return output
}

type ExtendedConnector struct {
	Connector *meroxa.Connector
	Logs      string
}

func extendedResourcesTable(resources []*meroxa.Resource, connectors []*meroxa.Connector) string {
	if len(resources) == 0 {
		return ""
	}
	subTable := "\tResources\n"

	var r *meroxa.Resource
	for _, c := range connectors {
		found := false
		for _, resource := range resources {
			if resource.Name == c.ResourceName {
				r = resource
				found = true
				break
			}
		}
		if !found {
			panic("internal error")
		}

		subTable += fmt.Sprintf("\t    %s (%s)\n", r.Name, string(c.Type))
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", r.UUID)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "Type", string(r.Type))
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "State", string(c.State))
	}

	return subTable
}

func extendedFunctionsTable(functions []*meroxa.Function) string {
	if len(functions) == 0 {
		return ""
	}
	subTable := "\tFunctions\n"

	for _, f := range functions {
		subTable += fmt.Sprintf("\t    %s\n", f.Name)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "UUID", f.UUID)
		subTable += fmt.Sprintf("\t\t%5s:   %s\n", "State", f.Status.State)
	}

	return subTable
}

func AppLogsTable(resources []meroxa.ApplicationResource, connectors []*ExtendedConnector, functions []*meroxa.Function) string {
	subTable := "\tResources:\n"

	var r meroxa.ApplicationResource
	for _, c := range connectors {
		found := false
		for _, resource := range resources {
			if resource.Name.String == c.Connector.ResourceName {
				r = resource
				found = true
				break
			}
		}
		if !found {
			panic("internal error")
		}

		if r.Collection.Source.String != "" {
			subTable += fmt.Sprintf("%s (source: %v)\n", r.Name.String, r.Collection.Source)
		}
		if r.Collection.Destination.String != "" {
			subTable += fmt.Sprintf("%s (destination: %v)\n", r.Name.String, r.Collection.Destination)
		}

		subTable += fmt.Sprintf("%s:\n%s\n", "Trace", c.Connector.Trace)
		subTable += fmt.Sprintf("%s:\n%s\n", "Logs", c.Logs)
	}

	subTable += "\tFunctions:\n"

	for _, f := range functions {
		subTable += fmt.Sprintf("%s Logs:\n%s\n", f.Name, f.Logs)
	}

	return subTable
}

// BuildsTable displays multiple build records for future listing per app.
func BuildsTable(builds []*meroxa.Build, hideHeaders bool) string {
	if len(builds) == 0 {
		return ""
	}

	table := simpletable.New()
	if !hideHeaders {
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "UUID"},
				{Align: simpletable.AlignCenter, Text: "STATE"},
				{Align: simpletable.AlignCenter, Text: "CREATED AT"},
				{Align: simpletable.AlignCenter, Text: "UPDATED AT"},
			},
		}
	}

	for _, p := range builds {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: p.Uuid},
			{Align: simpletable.AlignCenter, Text: p.Status.State},
			{Align: simpletable.AlignCenter, Text: p.CreatedAt},
			{Align: simpletable.AlignCenter, Text: p.UpdatedAt},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}

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

func truncateString(oldString string, l int) string {
	str := oldString

	if len(oldString) > l {
		str = oldString[:l] + "..."
	}

	return str
}

func PrettyString(a interface{}) (string, error) {
	j, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return "", err
	}
	if string(j) == "null" {
		return "", fmt.Errorf("unsuccessful marshal indent")
	}
	return string(j), nil
}
