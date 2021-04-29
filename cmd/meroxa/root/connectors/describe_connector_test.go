package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	mock "github.com/meroxa/cli/mock-cmd"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func TestDescribeConnectorArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires connector name"), name: ""},
		{args: []string{"connectorName"}, err: nil, name: "connectorName"},
	}

	for _, tt := range tests {
		ar := &DescribeConnector{}
		err := ar.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.args.Name)
		}
	}
}

func TestDescribeConnectorExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockDescribeConnectorClient(ctrl)
	logger := log.NewTestLogger()

	connectorName := "my-connector"

	c := utils.GenerateConnector(0, connectorName)
	client.
		EXPECT().
		GetConnectorByName(
			ctx,
			c.Name,
		).
		Return(&c, nil)

	dc := &DescribeConnector{
		client: client,
		logger: logger,
	}
	dc.args.Name = c.Name

	err := dc.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.ConnectorsTable([]*meroxa.Connector{&c})

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotConnector meroxa.Connector
	err = json.Unmarshal([]byte(gotJSONOutput), &gotConnector)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotConnector, c) {
		t.Fatalf("expected \"%v\", got \"%v\"", c, gotConnector)
	}
}
