package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestResourcesTable(t *testing.T) {
	resource := &meroxa.Resource{
		UUID:        "1dc8c9c6-d1d3-4b41-8f16-08302e87fc7b",
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "error",
		},
	}
	resIDAlign := &meroxa.Resource{
		UUID:        "9483768f-c384-4b4a-96bf-b80a79a23b5c",
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "ready",
		},
	}

	tests := map[string][]*meroxa.Resource{
		"Base":         {resource},
		"ID_Alignment": {resource, resIDAlign},
	}

	tableHeaders := []string{"ID", "NAME", "TYPE", "ENVIRONMENT", "URL", "TUNNEL", "STATE"}

	for name, resources := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintResourcesTable(resources, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			switch name {
			case "Base":
				if !strings.Contains(out, resource.Name) {
					t.Errorf("%s, not found", resource.Name)
				}
				if !strings.Contains(out, resource.UUID) {
					t.Errorf("%s, not found", resource.UUID)
				}
				if !strings.Contains(out, string(resource.Status.State)) {
					t.Errorf("state %s, not found", resource.Status.State)
				}
			case "ID_Alignment":
				if !strings.Contains(out, resIDAlign.Name) {
					t.Errorf("%s, not found", resIDAlign.Name)
				}
				if !strings.Contains(out, resIDAlign.UUID) {
					t.Errorf("%s, not found", resIDAlign.UUID)
				}
				if !strings.Contains(out, string(resIDAlign.Status.State)) {
					t.Errorf("state %s, not found", resource.Status.State)
				}
			}
			fmt.Println(out)
		})
	}
}

func TestResourcesTableWithoutHeaders(t *testing.T) {
	resource := &meroxa.Resource{
		UUID:        "9483768f-c384-4b4a-96bf-b80a79a23b5c",
		Type:        "jdbc",
		Name:        "my-db-jdbc-source",
		URL:         "postgres://display.test.us-east-1.rds.amazonaws.com:5432/display",
		Credentials: nil,
		Metadata:    nil,
		Status: meroxa.ResourceStatus{
			State: "error",
		},
	}

	var resources []*meroxa.Resource
	resources = append(resources, resource)

	tableHeaders := []string{"ID", "NAME", "TYPE", "URL", "TUNNEL", "STATE"}

	out := utils.CaptureOutput(func() {
		PrintResourcesTable(resources, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, resource.Name) {
		t.Errorf("%s, not found", resource.Name)
	}
	if !strings.Contains(out, resource.UUID) {
		t.Errorf("%s, not found", resource.UUID)
	}
	if !strings.Contains(out, string(resource.Status.State)) {
		t.Errorf("state %s, not found", resource.Status.State)
	}
}

var types = []meroxa.ResourceType{
	{
		Name:         string(meroxa.ResourceTypePostgres),
		ReleaseStage: meroxa.ResourceTypeReleaseStageBeta,
		FormConfig: map[string]interface{}{
			meroxa.ResourceTypeFormConfigHumanReadableKey: "PostgreSQL",
		},
	},
}

func TestResourceTypesTable(t *testing.T) {
	out := utils.CaptureOutput(func() {
		PrintResourceTypesTable(types, false)
	})

	if !strings.Contains(out, "NAME") {
		t.Errorf("NAME table headers is missing")
	}
	if !strings.Contains(out, "TYPE") {
		t.Errorf("TYPE table headers is missing")
	}
	if !strings.Contains(out, "RELEASE STAGE") {
		t.Errorf("table headers is missing")
	}

	for _, rType := range types {
		if !strings.Contains(
			out,
			fmt.Sprintf("%s", rType.FormConfig[meroxa.ResourceTypeFormConfigHumanReadableKey])) {
			t.Errorf("%s, not found", rType.Name)
		}
	}
}

func TestResourceTypesTableWithoutHeaders(t *testing.T) {
	out := utils.CaptureOutput(func() {
		PrintResourceTypesTable(types, true)
	})

	if strings.Contains(out, "NAME") {
		t.Errorf("NAME table headers unexpected")
	}
	if strings.Contains(out, "TYPE") {
		t.Errorf("TYPE table headers unexpected")
	}
	if strings.Contains(out, "RELEASE STAGE") {
		t.Errorf("RELEASE STAGE table headers unexpected")
	}

	for _, rType := range types {
		if !strings.Contains(
			out,
			fmt.Sprintf("%s", rType.FormConfig[meroxa.ResourceTypeFormConfigHumanReadableKey])) {
			t.Errorf("%s, not found", rType.Name)
		}
	}
}
