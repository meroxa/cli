package cmd

import (
	"bytes"
	"fmt"
	"github.com/meroxa/cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
	"testing"
)

func TestRemoveFlags(t *testing.T) {
	expectedFlags := []struct {
		name       string
		required   bool
		shorthand  string
		persistent bool
	}{
		{"force", false, "f", true},
	}

	c := &cobra.Command{}
	r := &Remove{}
	r.setFlags(c)

	for _, f := range expectedFlags {
		var cf *pflag.Flag

		if f.persistent {
			cf = c.PersistentFlags().Lookup(f.name)
		} else {
			cf = c.Flags().Lookup(f.name)
		}

		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
}

func TestConfirmationPrompt(t *testing.T) {
	expectedValue := "correct-value"

	tests := []struct {
		value string
	}{
		{expectedValue},
		{"incorrect-value"},
	}

	// Force flag to false
	r := Remove{false}

	for _, tt := range tests {
		output := utils.CaptureOutput(func() {
			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", expectedValue)))
			confirmed := r.confirmRemove(&stdin, tt.value)

			if confirmed && !strings.Contains(expectedValue, tt.value) {
				t.Fatalf("for value \"%s\", it shouldn't have been confirmed", tt.value)
			}

			if !confirmed && strings.Contains(expectedValue, tt.value) {
				t.Fatalf("for value \"%s\", it should have been confirmed", tt.value)
			}
		})

		expected := fmt.Sprintf("To proceed, type %s or re-run this command with --force\n", tt.value)

		if !strings.Contains(output, expected) {
			t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
		}
	}
}
