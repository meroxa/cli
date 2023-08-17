package display

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestBuildsLogsTable(t *testing.T) {
	build := "build"
	log := "custom log"
	now := time.Now().UTC()
	logs := meroxa.Logs{
		Data: []meroxa.LogData{
			{
				Timestamp: now,
				Log:       log,
				Source:    build,
			},
		},
		Metadata: meroxa.Metadata{
			End:   now,
			Start: now.Add(-12 * time.Hour),
			Limit: 10,
		},
	}

	out := BuildsLogsTable(&logs)

	if want := fmt.Sprintf("[%s]\t%q", now.Format(time.RFC3339), log); !strings.Contains(out, want) {
		t.Errorf("expected %q to be shown with logs, %s", want, out)
	}
}
