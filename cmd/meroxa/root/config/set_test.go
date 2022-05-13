package config

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"

	"testing"
)

var (
	errRequiredMinArgs          = errors.New("requires at least one KEY=VALUE pair (example: meroxa config set KEY=VALUE)")
	errRequiredBothKeyAndValue  = errors.New("there has to be at least a key before the `=` sign")
	errRequiredOnlyOneEqualSign = errors.New("a key=value needs to contain at least and only one `=` sign")
)

func TestSetConfigArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		keys map[string]string
	}{
		{args: nil, err: errRequiredMinArgs, keys: nil},
		{args: []string{"key=value"}, err: nil, keys: map[string]string{
			"KEY": "value",
		}},
		{args: []string{"key=value", "key2=value2"}, err: nil, keys: map[string]string{
			"KEY":  "value",
			"KEY2": "value2",
		}},
		{args: []string{"myGreatKey=value", "MY_EVEN_BETTER_KEY=value2"}, err: nil, keys: map[string]string{
			"MY_GREAT_KEY":       "value",
			"MY_EVEN_BETTER_KEY": "value2",
		}},
		{args: []string{"key="}, err: nil, keys: map[string]string{
			"KEY": "",
		}},
		{args: []string{"=value"}, err: errRequiredBothKeyAndValue, keys: nil},
		{args: []string{"key==value"}, err: errRequiredOnlyOneEqualSign, keys: nil},
		{args: []string{"key=value=value2", "MY_EVEN_BETTER_KEY=value2"}, err: errRequiredOnlyOneEqualSign, keys: nil},
	}

	for _, tt := range tests {
		s := &Set{}
		err := s.ParseArgs(tt.args)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected %q got %q", tt.err, err)
		}

		if err == nil && !reflect.DeepEqual(tt.keys, s.args.keys) {
			t.Fatalf("expected \"%v\" got \"%v\"", tt.keys, s.args.keys)
		}
	}
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{key: "key", expected: "KEY"},
		{key: "myKey", expected: "MY_KEY"},
		{key: "my_key", expected: "MY_KEY"},
		{key: "myAnotherKey", expected: "MY_ANOTHER_KEY"},
		{key: "KEY", expected: "KEY"},
		{key: "1KEY", expected: "1_KEY"},
		{key: "MY_KEY", expected: "MY_KEY"},
	}

	for _, tt := range tests {
		s := &Set{}
		got := s.normalizeKey(tt.key)

		if got != tt.expected {
			t.Fatalf("expected %q got %q", tt.expected, got)
		}
	}
}

func TestValidateAndAssignKeyValue(t *testing.T) {
	tests := []struct {
		keyValue      string
		err           error
		expectedKey   string
		expectedValue string
	}{
		{keyValue: "key==value", err: errRequiredOnlyOneEqualSign, expectedKey: "", expectedValue: ""},
		{keyValue: "key=", err: nil, expectedKey: "KEY", expectedValue: ""},
		{keyValue: "=value", err: errRequiredBothKeyAndValue, expectedKey: "", expectedValue: ""},
		{keyValue: "key=value", err: nil, expectedKey: "KEY", expectedValue: "value"},
		{keyValue: "KEY=value", err: nil, expectedKey: "KEY", expectedValue: "value"},
		{keyValue: "myKey=value", err: nil, expectedKey: "MY_KEY", expectedValue: "value"},
	}

	for _, tt := range tests {
		s := &Set{}
		s.args.keys = make(map[string]string)
		err := s.validateAndAssignKeyValue(tt.keyValue)

		if err != nil && tt.err.Error() != err.Error() {
			t.Fatalf("expected %q got %q", tt.err, err)
		}

		if err == nil && s.args.keys[tt.expectedKey] != tt.expectedValue {
			t.Fatalf("expected value %q for key %q, got %q", tt.expectedValue, tt.expectedKey, s.args.keys[tt.expectedKey])
		}
	}
}

func TestSetConfigExecution(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	k1 := "MY_KEY"
	k2 := "ANOTHER_KEY"

	setKeys := map[string]string{
		k1: "value",
		k2: "anotherValue",
	}

	cfg := config.NewInMemoryConfig()
	s := &Set{
		logger: logger,
		config: cfg,
	}

	s.args.keys = setKeys

	err := s.Execute(ctx)

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := fmt.Sprintf("Updating your Meroxa configuration file with %s=%s...\n"+
		"Updating your Meroxa configuration file with %s=%s...\n"+
		"Done!", k1, setKeys[k1], k2, setKeys[k2])

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	keys := []string{k1, k2}

	for _, k := range keys {
		if s.config.GetString(k) != setKeys[k] {
			t.Fatalf("expected value for key %q to be %q, got %q", k, setKeys[k], s.config.GetString(k))
		}
	}
}
