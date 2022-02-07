package utils

import (
	"bytes"
	"io"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/pflag"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func GeneratePipeline() meroxa.Pipeline {
	return meroxa.Pipeline{
		ID:    1,
		Name:  "pipeline-name",
		State: "healthy",
	}
}

func GeneratePipelineWithEnvironment() meroxa.Pipeline {
	p := GeneratePipeline()

	p.Environment = &meroxa.EnvironmentIdentifier{
		UUID: "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Name: "my-env",
	}

	return p
}

func GenerateResource() meroxa.Resource {
	return meroxa.Resource{
		ID:       1,
		Type:     meroxa.ResourceTypePostgres,
		Name:     "resource-1234",
		URL:      "https://user:password",
		Metadata: nil,
	}
}

func GenerateResourceWithEnvironment() meroxa.Resource {
	r := GenerateResource()

	r.Environment = &meroxa.EnvironmentIdentifier{
		UUID: "424ec647-9f0f-45a5-8e4b-3e0441f12444",
		Name: "my-environment",
	}
	return r
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
		Type:       meroxa.ConnectorTypeSource,
		Name:       connectorName,
		State:      meroxa.ConnectorStateRunning,
		PipelineID: pipelineID,
		Streams: map[string]interface{}{
			"output": []interface{}{"my-resource.Table"},
		},
	}
}

func GenerateConnectorWithEnvironment(pipelineID int, connectorName, envNameOrUUID string) meroxa.Connector {
	if pipelineID == 0 {
		pipelineID = rand.Intn(10000)
	}

	if connectorName == "" {
		connectorName = "connector-1234"
	}

	var env meroxa.EnvironmentIdentifier
	_, err := uuid.Parse(envNameOrUUID)
	if err == nil {
		env.UUID = envNameOrUUID
	} else {
		env.Name = envNameOrUUID
	}

	return meroxa.Connector{
		ID:         1,
		Type:       meroxa.ConnectorTypeSource,
		Name:       connectorName,
		State:      meroxa.ConnectorStateRunning,
		PipelineID: pipelineID,
		Streams: map[string]interface{}{
			"output": []interface{}{"my-resource.Table"},
		},
		Environment: &env,
	}
}

func GenerateEndpoint() meroxa.Endpoint {
	return meroxa.Endpoint{
		Name:              "endpoint",
		Protocol:          meroxa.EndpointProtocolHttp,
		Host:              "https://endpoint.test",
		Stream:            "stream",
		Ready:             true,
		BasicAuthUsername: "root",
		BasicAuthPassword: "secret",
	}
}

func GenerateTransform() meroxa.Transform {
	return meroxa.Transform{
		ID:          0,
		Name:        "ReplaceField",
		Required:    false,
		Description: "Filter or rename fields",
		Type:        "builtin",
		Properties:  nil,
	}
}

func GenerateEnvironment(environmentName string) meroxa.Environment {
	if environmentName == "" {
		environmentName = "environment-1234"
	}

	return meroxa.Environment{
		UUID:     "fd572375-77ce-4448-a071-ee4707a599d6",
		Type:     meroxa.EnvironmentTypePrivate,
		Name:     environmentName,
		Region:   meroxa.EnvironmentRegionUsEast1,
		Provider: meroxa.EnvironmentProviderAws,
		Status: meroxa.EnvironmentViewStatus{
			State: meroxa.EnvironmentStateProvisioned,
		},
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
