package display

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func TestLogsTable(t *testing.T) {
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

	out := LogsTable(&logs)

	if want := fmt.Sprintf("[%s]\t%s\t%q", now.Format(time.RFC3339), res, log); !strings.Contains(out, want) {
		t.Errorf("expected %q to be shown with logs, %s", want, out)
	}
	if want := fmt.Sprintf("[%s]\t%s\t%q", now.Format(time.RFC3339), fun, log); !strings.Contains(out, want) {
		t.Errorf("expected %q to be shown with logs, %s", want, out)
	}
}
