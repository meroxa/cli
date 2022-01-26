package functions

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreate(t *testing.T) {
	c := &Create{}

	m, err := c.parseEnvVars([]string{"KEY=VAL"})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(map[string]string{"KEY": "VAL"}, m); diff != "" {
		t.Fatalf("mismatch of parsed env vars (-want +got): %s", diff)
	}

	_, err = c.parseEnvVars([]string{"=VAL"})
	if err == nil {
		t.Fatal("error should not be nil")
	}
}
