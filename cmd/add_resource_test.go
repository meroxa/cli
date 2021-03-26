package cmd

import (
	"encoding/json"
	utils "github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"testing"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err error
		name string
	}{
		{[]string{""},nil, ""},
		{[]string{"resName"},nil, "resName"},
	}

	for _, tt := range tests {
		name, err := AddResource{}.checkArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, name)
		}
	}
}


func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name string
		required bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"credentials", false, ""},
		{"metadata", false, "m"},
	}

	c := AddResource{}.setFlags(&cobra.Command{})

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
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

func TestAddResourceOutput(t *testing.T) {
	r := utils.GenerateResource()

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	expected := "Resource resource-name successfully added!"

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestAddResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)


	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}

// TODO: Test adddResource


