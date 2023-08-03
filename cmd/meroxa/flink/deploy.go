package flink

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/ir"
	"golang.org/x/mod/semver"
)

const (
	irFilename       = "meroxa-ir.json"
	majorJavaVersion = "v11"
	modeEnvVar       = "MEROXA_PLATFORM"
	outputEnvVar     = "MEROXA_OUTPUT"
	irVal            = "EMIT_IR"
)

func doesJobUseMeroxaPlatform(ctx context.Context, jarPath string) (bool, error) {
	//TODO: This approach is fine for now but may prove to be too naive in the future
	// https://github.com/meroxa/product/issues/953
	cmd := exec.CommandContext(ctx, "jar", "-tf", jarPath)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "com.meroxa") {
			return true, nil
		}
	}

	return false, nil
}

func GetIRSpec(ctx context.Context, jarPath string, secrets map[string]string, l log.Logger) (*ir.DeploymentSpec, error) {
	if os.Getenv("UNIT_TEST") != "" {
		return nil, nil
	}

	verifyJavaVersion(ctx, l)

	usesMeroxa, err := doesJobUseMeroxaPlatform(ctx, jarPath)
	if err != nil {
		return nil, err
	}

	if !usesMeroxa {
		return nil, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	irFilepath := filepath.Join(cwd, irFilename)

	// The submitted jar is executed with some special env vars set to inform the `MeroxaExecutionEnvironment` to
	// short circuit execution and emit an IR spec instead
	// https://github.com/meroxa/flink-platform-prototype/blob/main/src/main/java/com/meroxa/flink/MeroxaExecutionEnvironment.java#L64-L69
	cmd := exec.CommandContext(ctx, "java", "-jar", jarPath)
	cmd.Env = append(
		cmd.Environ(),
		fmt.Sprintf("%s=%s", modeEnvVar, irVal),
		fmt.Sprintf("%s=%s", outputEnvVar, irFilename))
	output, err := cmd.CombinedOutput() // all java output goes to stderr, so that's fun
	defer func() {
		_ = os.Remove(irFilepath)
	}()
	if err != nil {
		return nil, fmt.Errorf("%s\n%v", output, err)
	}

	_, err = os.Stat(irFilepath)
	if err != nil {
		// @TODO try the docker way because the jar is skinny
		// Otherwise, there are no Meroxa* classes in this main class
		return nil, nil
	}

	b, err := os.ReadFile(irFilepath)
	if err != nil {
		return nil, err
	}

	// @TODO assess the scope of updating validateCollections to use the ConnectorSpec
	var spec ir.DeploymentSpec
	if err = json.Unmarshal(b, &spec); err != nil {
		return nil, err
	}

	spec.Secrets = secrets
	// workarounds to validate spec
	spec.Definition.Metadata.SpecVersion = ir.LatestSpecVersion
	spec.Definition.Metadata.Turbine.Language = "js" // java and flink are not acceptable yet
	spec.Definition.Metadata.Turbine.Version = majorJavaVersion

	// hardcode all sources to one destination as streams
	destinationUUID := ""
	var sourceUUIDs []string
	for _, cs := range spec.Connectors {
		if cs.Type == ir.ConnectorDestination {
			destinationUUID = cs.UUID
		} else {
			sourceUUIDs = append(sourceUUIDs, cs.UUID)
		}
	}

	for _, u := range sourceUUIDs {
		ss := ir.StreamSpec{
			UUID:     uuid.New().String(),
			FromUUID: u,
			ToUUID:   destinationUUID,
			Name:     u + "_" + destinationUUID,
		}
		spec.Streams = append(spec.Streams, ss)
	}
	return &spec, nil
}

func verifyJavaVersion(ctx context.Context, l log.Logger) {
	cmd := exec.CommandContext(ctx, "java", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.Warnf(ctx,
			"warning: unable to verify local Java version is compatible with the Meroxa Platform; jar's must be compiled for %s",
			majorJavaVersion)
		return
	}

	// looks like 'openjdk version "11.0.19" 2023-04-18'
	r := regexp.MustCompile(`version "([0-9.]+.[0-9.]+.[0-9.]+)"`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) > 0 {
		version := "v" + matches[1]
		comparison := semver.Compare(semver.Major(version), majorJavaVersion)
		if comparison == 0 {
			return
		}
		l.Warnf(ctx,
			"warning: local Java version %q is incompatible with the Meroxa Platform; jar's must be compiled for %s",
			version,
			majorJavaVersion)
	}
	return
}
