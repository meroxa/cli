package resources

import (
	"errors"
	"testing"
)

func TestRotateTunnelKeyArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{args: nil, err: errors.New("requires resource name"), name: ""},
		{args: []string{"resource-name"}, err: nil, name: "resource-name"},
	}

	for _, tt := range tests {
		cc := &RotateTunnelKey{}
		err := cc.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != cc.args.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, cc.args.Name)
		}
	}
}
