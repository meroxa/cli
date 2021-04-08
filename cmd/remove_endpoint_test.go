package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	"github.com/meroxa/cli/utils"
	"strings"
	"testing"
)

func TestRemoveEndpointArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, errors.New("requires endpoint name\n\nUsage:\n  meroxa remove endpoint NAME [flags]"), ""},
		{[]string{"endpoint-name"}, nil, "endpoint-name"},
	}

	r := &Remove{}
	for _, tt := range tests {
		rr := &RemoveEndpoint{removeCmd: r}
		err := rr.setArgs(tt.args)

		if tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != rr.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, rr.name)
		}

		if err == nil {
			componentType := "endpoint"
			if rr.removeCmd.componentType != componentType {
				t.Fatalf("expected type to be set to %q", componentType)
			}

			if rr.removeCmd.confirmableName != rr.name {
				t.Fatalf("expected \"confirmableName\" to be set to %q", rr.name)
			}
		}
	}
}

func TestRemoveEndpointExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockRemoveEndpointClient(ctrl)

	endpointName := "my-endpoint"

	client.
		EXPECT().
		DeleteEndpoint(ctx, endpointName).
		Return(nil)

	rc := &Remove{}

	re := &RemoveEndpoint{
		name:      endpointName,
		removeCmd: rc,
	}
	err := re.execute(ctx, client)

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestRemoveEndpointOutput(t *testing.T) {
	re := &RemoveEndpoint{name: "endpoint-name"}

	output := utils.CaptureOutput(func() {
		re.output()
	})

	expected := fmt.Sprintf("endpoint %s successfully removed", re.name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}
