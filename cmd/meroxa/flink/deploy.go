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

const irFilename = "meroxa-ir.json"
const exactJavaVersion = "v11"
const modeEnvVar = "MEROXA_PLATFORM"
const outputEnvVar = "MEROXA_OUTPUT"
const irVal = "EMIT_IR"

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
	_, err = cmd.CombinedOutput() // all java output goes to stderr
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

	// when integrating this with Turbine, convert to turbine.ApplicationResource format to run validateCollections
	// or assess the scope of updating validateCollections tp use the ConnectorSpec... aka breaking backwards compatibility
	var spec ir.DeploymentSpec
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}

	spec.Secrets = secrets
	return &spec, nil
}

func verifyJavaVersion(ctx context.Context, l log.Logger) {
	/*
			// more reliable way with bytecode
		type class struct {
			Magic        uint32
			MinorVersion uint16
			MajorVersion uint16
		}

		var versions = map[uint16]string{
			52: "Java 8",
			53: "Java 9",
			54: "Java 10",
			55: "Java 11",
			56: "Java 12",
			57: "Java 13",
			58: "Java 14",
			59: "Java 15",
			60: "Java 16",
			61: "Java 17",
		}
			unzip -p target/meroxa-flink-prototype-0.1.0-SNAPSHOT.jar META-INF/MANIFEST.MF
			"Main-Class" from that output
		filename = cwd + mainClass // we may not be at the top of the job definition or it may not exist on this system at all
				// cl = Readfile(filename)
				b := make([]byte, binary.Size(cl))
				f, err := os.Open(os.Args[1])
				if err != nil {
					log.Fatal(err)
				}
				if _, err := f.Read(b); err != nil {
					log.Fatal(err)
				}
				buf := bytes.NewReader(b)
				if err := binary.Read(buf, binary.BigEndian, &cl); err != nil {
					log.Fatal(err)
				}
				log.Printf("version: %s\n", versions[cl.MajorVersion])
	*/
	cmd := exec.CommandContext(ctx, "java", "-version")
	output, err := cmd.CombinedOutput() // everything goes to stderr
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
