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
		value           string
		confirmableName string
		componentType   string
	}{
		{expectedValue, expectedValue, "resource"},
		{"incorrect-value", expectedValue, "connector"},
		{"incorrect-value", expectedValue, ""},
	}

	// Force flag to false
	r := &Remove{}

	for _, tt := range tests {
		var err error
		r.confirmableName = tt.confirmableName
		r.componentType = tt.componentType

		output := utils.CaptureOutput(func() {
			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", expectedValue)))
			err = r.confirmRemove(&stdin, tt.value)

			if err == nil && !strings.Contains(expectedValue, tt.value) {
				t.Fatalf("for value \"%s\", it shouldn't have been confirmed", tt.value)
			}

			if err != nil && strings.Contains(expectedValue, tt.value) {
				t.Fatalf("for value \"%s\", it should have been confirmed", tt.value)
			}
		})

		expected := fmt.Sprintf("To proceed, type %s or re-run this command with --force\n", tt.value)

		if !strings.Contains(output, expected) {
			t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
		}

		if err != nil {
			if r.componentType != "" {
				expected = fmt.Sprintf("removing %s not confirmed", r.componentType)
			} else {
				expected = "removing value not confirmed"
			}

			if err.Error() != expected {
				t.Fatalf("not expected error, got \"%s\"", err.Error())
			}
		}
	}
}

func TestAddConfirmation(t *testing.T) {
	r := &Remove{}
	cmd := r.command()
	r.force = true

	r.confirmableName = "argument-name"

	for _, c := range cmd.Commands() {
		output := utils.CaptureOutput(func() {
			err := c.PreRunE(c, []string{r.confirmableName})

			if err != nil {
				t.Fatalf("not expected error, got \"%s\"", err.Error())
			}
		})

		expected := fmt.Sprintf("Removing %s...\n", r.confirmableName)

		if !strings.Contains(expected, output) {
			t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
		}
	}
}
