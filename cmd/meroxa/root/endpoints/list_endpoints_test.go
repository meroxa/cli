package endpoints

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func getEndpoints() []meroxa.Endpoint {
	var endpoints []meroxa.Endpoint
	e := utils.GenerateEndpoint()
	return append(endpoints, e)
}

func TestListConnectorsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMocklistEndpointsClient(ctrl)
	logger := log.NewTestLogger()

	endpoints := getEndpoints()

	client.
		EXPECT().
		ListEndpoints(ctx).
		Return(endpoints, nil)

	l := &ListEndpoints{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := utils.EndpointsTable(endpoints)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotEndpoints []meroxa.Endpoint
	err = json.Unmarshal([]byte(gotJSONOutput), &gotEndpoints)

	var le []meroxa.Endpoint

	le = append(le, endpoints...)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotEndpoints, le) {
		t.Fatalf("expected \"%v\", got \"%v\"", endpoints, gotEndpoints)
	}
}
