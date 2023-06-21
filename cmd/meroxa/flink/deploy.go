package flink

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/ir"

	"golang.org/x/mod/semver"
)

const (
	irFilename       = "meroxa-ir.json"
	exactJavaVersion = "v11"
	modeEnvVar       = "MEROXA_PLATFORM"
	outputEnvVar     = "MEROXA_OUTPUT"
	irVal            = "EMIT_IR"
)

func GetIRSpec(ctx context.Context, jarPath string, secrets map[string]string, l log.Logger) (*ir.DeploymentSpec, error) {
	verifyJavaVersion(ctx, l)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	irFilepath := filepath.Join(cwd, irFilename)

	cmd := exec.CommandContext(ctx, "java", "-jar", jarPath)
	cmd.Env = append(
		cmd.Environ(),
		fmt.Sprintf("%s=%s", modeEnvVar, irVal),
		fmt.Sprintf("%s=%s", outputEnvVar, irFilename))
	_, err = cmd.CombinedOutput() // all java output goes to stderr, so that's fun
	defer func() {
		_ = os.Remove(irFilepath)
	}()
	if err != nil {
		return nil, err
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
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}

	spec.Secrets = secrets
	return &spec, nil
}

func verifyJavaVersion(ctx context.Context, l log.Logger) {
	cmd := exec.CommandContext(ctx, "java", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		l.Warnf(ctx,
			"warning: unable to verify local Java version is compatible with the Meroxa Platform; jar's must be compiled for %s",
			exactJavaVersion)
		return
	}

	// looks like 'openjdk version "11.0.19" 2023-04-18'
	r := regexp.MustCompile(`version "([0-9.]+.[0-9.]+.[0-9.]+)"`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) > 0 {
		version := "v" + matches[1]
		comparison := semver.Compare(version, exactJavaVersion)
		if comparison >= 1 {
			return
		}
		l.Warnf(ctx,
			"warning: local Java version %q is incompatible with the Meroxa Platform; jar's must be compiled for %s",
			version,
			exactJavaVersion)
	}
	return
}
