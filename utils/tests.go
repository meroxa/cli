package utils

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/pflag"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func GeneratePipeline() meroxa.Pipeline {
	return meroxa.Pipeline{
		UUID:  "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Name:  "pipeline-name",
		State: "healthy",
	}
}

func GeneratePipelineWithEnvironment() meroxa.Pipeline {
	p := GeneratePipeline()

	p.Environment = &meroxa.EntityIdentifier{
		UUID: "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Name: "my-env",
	}

	return p
}

func GenerateResourceWithNameAndStatus(resourceName, resourceState string) meroxa.Resource {
	if resourceName == "" {
		resourceName = "resource-1234"
	}

	newResource := meroxa.Resource{
		UUID:     "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Type:     meroxa.ResourceTypePostgres,
		Name:     resourceName,
		URL:      "https://user:password",
		Metadata: nil,
	}

	if resourceState == "ready" {
		newResource.Status.State = meroxa.ResourceStateReady
	}

	return newResource
}

func GenerateResourceWithEnvironment() meroxa.Resource {
	r := GenerateResource()

	r.Environment = &meroxa.EntityIdentifier{
		UUID: "424ec647-9f0f-45a5-8e4b-3e0441f12444",
		Name: "my-environment",
	}
	return r
}

func GenerateResource() meroxa.Resource {
	return GenerateResourceWithNameAndStatus("", "")
}

func GenerateConnector(pipelineName, connectorName string) meroxa.Connector {
	if pipelineName == "" {
		pipelineName = "pipeline-1234"
	}

	if connectorName == "" {
		connectorName = "connector-1234"
	}

	return meroxa.Connector{
		UUID:         "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Type:         meroxa.ConnectorTypeSource,
		Name:         connectorName,
		State:        meroxa.ConnectorStateRunning,
		PipelineName: pipelineName,
		Streams: map[string]interface{}{
			"output": []interface{}{"my-resource.Table"},
		},
	}
}

func GenerateConnectorWithEnvironment(pipelineName, connectorName, envNameOrUUID string) meroxa.Connector {
	if pipelineName == "" {
		pipelineName = "pipeline-1234"
	}

	if connectorName == "" {
		connectorName = "connector-1234"
	}

	var env meroxa.EntityIdentifier
	_, err := uuid.Parse(envNameOrUUID)
	if err == nil {
		env.UUID = envNameOrUUID
	} else {
		env.Name = envNameOrUUID
	}

	return meroxa.Connector{
		UUID:         "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Type:         meroxa.ConnectorTypeSource,
		Name:         connectorName,
		State:        meroxa.ConnectorStateRunning,
		PipelineName: pipelineName,
		Streams: map[string]interface{}{
			"output": []interface{}{"my-resource.Table"},
		},
		Environment: &env,
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

func GenerateEnvironmentFailed(environmentName string) meroxa.Environment {
	if environmentName == "" {
		environmentName = "environment-1234-bad"
	}

	return meroxa.Environment{
		UUID:     "fd572375-77ce-4448-a071-ee4707a599d6",
		Type:     meroxa.EnvironmentTypePrivate,
		Name:     environmentName,
		Region:   meroxa.EnvironmentRegionUsEast1,
		Provider: meroxa.EnvironmentProviderAws,
		Status: meroxa.EnvironmentViewStatus{
			State:   meroxa.EnvironmentStatePreflightError,
			Details: "",
			PreflightDetails: &meroxa.PreflightDetails{
				PreflightPermissions: &meroxa.PreflightPermissions{
					S3:  []string{"missing read permission for S3", "missing write permissions for S3"},
					EC2: []string{"missing read permission for S3", "missing write permissions for S3"},
				},
				PreflightLimits: &meroxa.PreflightLimits{
					EIP: "",
				},
			},
		},
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
			State: meroxa.EnvironmentStatePreflightSuccess,
		},
	}
}

func GenerateApplication(state meroxa.ApplicationState) meroxa.Application {
	if state == "" {
		state = meroxa.ApplicationStateRunning
	}
	return meroxa.Application{
		Name:     "application-name",
		Language: "golang",
		Status: meroxa.ApplicationStatus{
			State: state,
		},
	}
}

func GenerateApplicationWithEnv(
	state meroxa.ApplicationState,
	envType meroxa.EnvironmentType,
	provider meroxa.EnvironmentProvider,
) meroxa.Application {
	app := GenerateApplication(state)
	app.Environment = meroxa.EntityIdentifier{
		UUID: uuid.NewString(), Name: "my-env",
	}
	return app
}

func GenerateBuild() meroxa.Build {
	return meroxa.Build{
		Uuid: "236d6e81-6a22-4805-b64f-3fa0a57fdbdc",
		Status: meroxa.BuildStatus{
			State:   "status",
			Details: "details",
		},
		CreatedAt:  time.Time{}.String(),
		UpdatedAt:  time.Time{}.String(),
		SourceBlob: meroxa.SourceBlob{Url: "url"},
		Image:      "image",
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
