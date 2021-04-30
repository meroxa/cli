package connectors

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
)

func TestLogsConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires connector name"), name: ""},
		{args: []string{"conName"}, err: nil, name: "conName"},
	}

	for _, tt := range tests {
		cc := &LogsConnector{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestLogsConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockLogsConnectorClient(ctrl)
	logger := log.NewTestLogger()

	connectorName := "connector-name"

	c := &LogsConnector{
		client: client,
		logger: logger,
	}

	c.args.Name = connectorName

	var responseDetails = ioutil.NopCloser(bytes.NewReader([]byte(
		`[2021-04-29T12:16:42Z] Just another log line from my connector`,
	)))

	var httpResponse = &http.Response{
		StatusCode: 200,
		Body:       responseDetails,
	}

	client.
		EXPECT().
		GetConnectorLogs(ctx, connectorName).
		Return(httpResponse, nil)

	err := c.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := "[2021-04-29T12:16:42Z] Just another log line from my connector"

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}
}
