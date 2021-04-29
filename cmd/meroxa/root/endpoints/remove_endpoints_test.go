package endpoints

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
)

func TestRemoveEndpointArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires endpoint name"), name: ""},
		{args: []string{"endpoint-name"}, err: nil, name: "endpoint-name"},
	}

	for _, tt := range tests {
		cc := &RemoveEndpoint{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}

func TestRemoveEndpointExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemoveEndpointClient(ctrl)
	logger := log.NewTestLogger()

	r := &RemoveEndpoint{
		client: client,
		logger: logger,
	}

	e := utils.GenerateEndpoint()
	r.args.Name = e.Name

	client.
		EXPECT().
		DeleteEndpoint(ctx, r.args.Name).
		Return(nil)

	err := r.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf(`Removing endpoint %q...
Endpoint %q successfully removed
`, r.args.Name, r.args.Name)

	if gotLeveledOutput != wantLeveledOutput {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}
}
