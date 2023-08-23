package flink

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/meroxa/cli/utils"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestLogsFlinkJobArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires Flink Job name or UUID"), name: ""},
		{args: []string{"job-name"}, err: nil, name: "job-name"},
	}

	for _, tt := range tests {
		l := &Logs{}
		err := l.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != l.args.NameOrUUID {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, l.args.NameOrUUID)
		}
	}
}

func TestFlinkLogsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	fj := utils.GenerateFlinkJob()

	flinkLogs := &meroxa.Logs{
		Data: []meroxa.LogData{
			{
				Timestamp: time.Now().UTC(),
				Log:       "log just logging",
				Source:    "connector",
			},
			{
				Timestamp: time.Now().UTC(),
				Log:       "another log",
				Source:    "flink-job",
			},
		},
		Metadata: meroxa.Metadata{
			End:   time.Now().UTC(),
			Start: time.Now().UTC().Add(-12 * time.Hour),
			Limit: 10,
		},
	}

	client.EXPECT().GetFlinkLogsV2(ctx, fj.Name).Return(flinkLogs, nil)

	l := &Logs{
		client: client,
		logger: logger,
	}
	l.args.NameOrUUID = fj.Name

	err := l.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.LogsTable(flinkLogs)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf(cmp.Diff(wantLeveledOutput, gotLeveledOutput))
	}

	gotJSONOutput := logger.JSONOutput()
	var gotAppLogs meroxa.Logs
	err = json.Unmarshal([]byte(gotJSONOutput), &gotAppLogs)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotAppLogs, *flinkLogs) {
		t.Fatalf(cmp.Diff(*flinkLogs, gotAppLogs))
	}
}
