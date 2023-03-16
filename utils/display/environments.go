package display

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

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

//nolint:funlen
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
				{Align: simpletable.AlignRight, Text: "AWS ECR Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.ECR, " ; ")},
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
				{Align: simpletable.AlignRight, Text: "AWS ECR Permissions Status:"},
				{Text: strings.Join(environment.Status.PreflightDetails.PreflightPermissions.ECR, " ; ")},
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
