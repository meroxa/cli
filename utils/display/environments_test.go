package display

import (
	"fmt"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestEnvironmentsTable(t *testing.T) {
	e := &meroxa.Environment{
		Type:     meroxa.EnvironmentTypePrivate,
		Name:     "environment-1234",
		Provider: meroxa.EnvironmentProviderAws,
		Region:   meroxa.EnvironmentRegionUsEast1,
		Status:   meroxa.EnvironmentViewStatus{State: meroxa.EnvironmentStateReady},
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
	}

	tests := map[string][]*meroxa.Environment{
		"Base": {e},
	}

	tableHeaders := []string{"ID", "NAME", "TYPE", "PROVIDER", "REGION", "STATE"}

	for name, environments := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintEnvironmentsTable(environments, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}
			if !strings.Contains(out, e.Name) {
				t.Errorf("%s, not found", e.Name)
			}
			if !strings.Contains(out, string(e.Type)) {
				t.Errorf("%s, not found", e.Type)
			}
			if !strings.Contains(out, string(e.Region)) {
				t.Errorf("%s, not found", e.Region)
			}
			if !strings.Contains(out, string(e.Status.State)) {
				t.Errorf("%s, not found", e.Status.State)
			}
			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}

			fmt.Println(out)
		})
	}
}

func TestEnvironmentsTablePreflightFailed(t *testing.T) {
	e := utils.GenerateEnvironmentFailed("environment-preflight-failed")

	tests := map[string][]*meroxa.Environment{
		"Base": {&e},
	}

	tableHeaders := []string{"ID", "NAME", "TYPE", "PROVIDER", "REGION", "STATE"}

	for name, environments := range tests {
		t.Run(name, func(t *testing.T) {
			out := utils.CaptureOutput(func() {
				PrintEnvironmentsTable(environments, false)
			})

			for _, header := range tableHeaders {
				if !strings.Contains(out, header) {
					t.Errorf("%s header is missing", header)
				}
			}

			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}
			if !strings.Contains(out, e.Name) {
				t.Errorf("%s, not found", e.Name)
			}
			if !strings.Contains(out, string(e.Type)) {
				t.Errorf("%s, not found", e.Type)
			}
			if !strings.Contains(out, string(e.Region)) {
				t.Errorf("%s, not found", e.Region)
			}
			if !strings.Contains(out, string(e.Status.State)) {
				t.Errorf("%s, not found", e.Status.State)
			}
			if !strings.Contains(out, e.UUID) {
				t.Errorf("%s, not found", e.UUID)
			}

			fmt.Println(out)
		})
	}
}

func TestEnvironmentsTableWithoutHeaders(t *testing.T) {
	e := &meroxa.Environment{
		Type:     meroxa.EnvironmentTypePrivate,
		Name:     "environment-1234",
		Provider: meroxa.EnvironmentProviderAws,
		Region:   meroxa.EnvironmentRegionUsEast1,
		Status:   meroxa.EnvironmentViewStatus{State: meroxa.EnvironmentStateReady},
		UUID:     "531428f7-4e86-4094-8514-d397d49026f7",
	}

	var environments []*meroxa.Environment
	tableHeaders := []string{"ID", "NAME", "TYPE", "PROVIDER", "REGION", "STATE"}

	environments = append(environments, e)

	out := utils.CaptureOutput(func() {
		PrintEnvironmentsTable(environments, true)
	})

	for _, header := range tableHeaders {
		if strings.Contains(out, header) {
			t.Errorf("%s header should not be displayed", header)
		}
	}

	if !strings.Contains(out, e.UUID) {
		t.Errorf("%s, not found", e.UUID)
	}
	if !strings.Contains(out, e.Name) {
		t.Errorf("%s, not found", e.Name)
	}
	if !strings.Contains(out, string(e.Type)) {
		t.Errorf("%s, not found", e.Type)
	}
	if !strings.Contains(out, string(e.Region)) {
		t.Errorf("%s, not found", e.Region)
	}
	if !strings.Contains(out, string(e.Status.State)) {
		t.Errorf("%s, not found", e.Status.State)
	}
	if !strings.Contains(out, e.UUID) {
		t.Errorf("%s, not found", e.UUID)
	}
}
