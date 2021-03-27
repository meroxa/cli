package utils

import (
	"bytes"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/pflag"
	"io"
	"os"
)

func GenerateResource() meroxa.Resource {
	return meroxa.Resource{
		ID: 1,
		Type: "postgres",
		Name: "resource-1234",
		URL: "https://user:password",
		Metadata: nil,
	}
}

func IsFlagRequired(flag *pflag.Flag) bool{
	requiredAnnotation := "cobra_annotation_bash_completion_one_required_flag"

	if len(flag.Annotations[requiredAnnotation]) > 0 && flag.Annotations[requiredAnnotation][0] == "true" {
		return true
	}

	return false
}

// CaptureOutput is used to capture stdout to be compared on tests
func CaptureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	print()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	f()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out
}
