package flink

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GetVersion will return the tag or hash of the turbine-go dependency of a given app.
func (t *flinkCLI) GetVersion(ctx context.Context) (string, error) {
	var cmd *exec.Cmd

	cmd = exec.CommandContext(
		ctx,
		"java",
		"--version")
	fmtErr := fmt.Errorf(
		"unable to determine the version of java used by the Meroxa Flink Job at %s",
		t.appPath)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmtErr
	}
	/*
		openjdk 11.0.18 2023-01-17
		OpenJDK Runtime Environment (build 11.0.18+10-post-Ubuntu-0ubuntu120.04.1)
		OpenJDK 64-Bit Server VM (build 11.0.18+10-post-Ubuntu-0ubuntu120.04.1, mixed mode, sharing)
	*/

	version := strings.TrimSpace(string(stdout))
	return version, nil
}
