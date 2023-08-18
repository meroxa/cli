package display

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/meroxa/cli/utils"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestAppLogsTable(t *testing.T) {
	res := "res"
	fun := "fun"
	log := "custom log"
	now := time.Now().UTC()
	logs := meroxa.Logs{
		Data: []meroxa.LogData{
			{
				Timestamp: now,
				Log:       log,
				Source:    fun,
			},
			{
				Timestamp: now,
				Log:       log,
				Source:    res,
			},
		},
		Metadata: meroxa.Metadata{
			End:   now,
			Start: now.Add(-12 * time.Hour),
			Limit: 10,
		},
	}

	out := AppLogsTable(&logs)

	if want := fmt.Sprintf("[%s]\t%s\t%q", now.Format(time.RFC3339), res, log); !strings.Contains(out, want) {
		t.Errorf("expected %q to be shown with logs, %s", want, out)
	}
	if want := fmt.Sprintf("[%s]\t%s\t%q", now.Format(time.RFC3339), fun, log); !strings.Contains(out, want) {
		t.Errorf("expected %q to be shown with logs, %s", want, out)
	}
}

func TestAppDescribeTable(t *testing.T) {
	testCases := []struct {
		desc                 string
		app                  func() *meroxa.Application
		shouldErrorOnEnvInfo func(string) bool
	}{
		{
			desc: "Application with no environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplication("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return strings.Contains(output, "Environment")
			},
		},
		{
			desc: "Application with in a private environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplicationWithEnv("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, "Environment")
			},
		},
		{
			desc: "Application with in a private environment",
			app: func() *meroxa.Application {
				a := utils.GenerateApplicationWithEnv("")
				return &a
			},
			shouldErrorOnEnvInfo: func(output string) bool {
				return !strings.Contains(output, "Environment")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			app := tc.app()
			out := AppTable(app)
			if !strings.Contains(out, "UUID") {
				t.Errorf("expected %q to be shown\n%s\n", app.UUID, out)
			}
			if !strings.Contains(out, "Name:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Name, out)
			}
			if !strings.Contains(out, "Language:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Language, out)
			}
			if !strings.Contains(out, "Git SHA:") {
				t.Errorf("expected %q to be shown\n%s\n", app.GitSha, out)
			}
			if !strings.Contains(out, "Created At:") {
				t.Errorf("expected %q to be shown\n%s\n", app.CreatedAt, out)
			}
			if !strings.Contains(out, "Updated At:") {
				t.Errorf("expected %q to be shown\n%s\n", app.UpdatedAt, out)
			}
			if !strings.Contains(out, "State:") {
				t.Errorf("expected %q to be shown\n%s\n", app.Status.State, out)
			}
			if tc.shouldErrorOnEnvInfo(out) {
				t.Errorf("expected environment information to be shown\n%s\n", out)
			}
		})
	}
}
