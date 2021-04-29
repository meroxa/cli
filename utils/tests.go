package utils

import (
	"bytes"
	"io"
	"math/rand"
	"os"

	"github.com/meroxa/meroxa-go"
	"github.com/spf13/pflag"
)

func GeneratePipeline() meroxa.Pipeline {
	return meroxa.Pipeline{
		ID:       1,
		Name:     "pipeline-name",
		Metadata: nil,
		State:    "healthy",
	}
}

func GenerateResource() meroxa.Resource {
	return meroxa.Resource{
		ID:       1,
		Type:     "postgres",
		Name:     "resource-1234",
		URL:      "https://user:password",
		Metadata: nil,
	}
}

func GenerateConnector(pipelineID int, connectorName string) meroxa.Connector {
	if pipelineID == 0 {
		pipelineID = rand.Intn(10000)
	}

	if connectorName == "" {
		connectorName = "connector-1234"
	}

	return meroxa.Connector{
		ID:         1,
		Type:       "postgres",
		Name:       connectorName,
		State:      "running",
		PipelineID: pipelineID,
		Streams: map[string]interface{}{
			"output": []interface{}{"my-resource.Table"},
		},
	}
}

func GenerateEndpoint() meroxa.Endpoint {
	return meroxa.Endpoint{
		Name:              "endpoint",
		Protocol:          "http",
		Host:              "https://endpoint.test",
		Stream:            "stream",
		Ready:             true,
		BasicAuthUsername: "root",
		BasicAuthPassword: "secret",
	}
}

func IsFlagRequired(flag *pflag.Flag) bool {
	requiredAnnotation := "cobra_annotation_bash_completion_one_required_flag"

	if len(flag.Annotations[requiredAnnotation]) > 0 && flag.Annotations[requiredAnnotation][0] == "true" {
		return true
	}

	return false
}

// CaptureOutput is used to capture stdout to be compared on tests.
func CaptureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	print()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()
	f()

	// back to normal state
	_ = w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	return out
}
