package display

import (
	"strings"
	"testing"
)

type genericType struct {
	ABC string `json:"abc"`
	DEF string `json:"def"`
}

func TestGenericPrints(t *testing.T) {
	forTable := genericType{ABC: "abd", DEF: "def"}
	details := map[string]string{"Abc": "abc", "Def": "def"}
	out := PrintTable(forTable, details)
	for key := range details {
		if !strings.Contains(out, key) {
			t.Errorf("Output does not contain the label %q", key)
		}
	}

	forList := []genericType{forTable, forTable}
	details = map[string]string{"ABC": "abc", "DED": "def"}
	out = PrintList(forList, details)
	for key := range details {
		if !strings.Contains(out, key) {
			t.Errorf("Output does not contain the header %q", key)
		}
	}
}
