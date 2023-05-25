package turbine

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/meroxa/cli/log"
)

type Git struct{}

func (g *Git) GitInit(ctx context.Context, appPath string) error {
	if appPath == "" {
		return errors.New("path is required")
	}

	isGitOlderThan228, err := checkGitVersion(ctx)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "git", "init", appPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(string(out))
	}

	if !isGitOlderThan228 {
		cmd = exec.CommandContext(ctx, "git", "config", "init.defaultBranch", "main")
		cmd.Dir = appPath
		if out, err := cmd.CombinedOutput(); err != nil {
			return errors.New(string(out))
		}
	}
	cmd = exec.CommandContext(ctx, "git", "checkout", "-b", "main")
	cmd.Dir = appPath
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.New(string(out))
	}

	return nil
}

// CheckUncommittedChanges prints warnings about uncommitted tracked and untracked files.
func (g *Git) CheckUncommittedChanges(ctx context.Context, l log.Logger, appPath string) error {
	l.Info(ctx, "Checking for uncommitted changes...")
	cmd := exec.Command("git", "status", "--porcelain=v2")
	cmd.Dir = appPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	all := string(output)
	lines := strings.Split(strings.TrimSpace(all), "\n")
	if len(lines) > 0 && lines[0] != "" {
		cmd = exec.Command("git", "status")
		output, err = cmd.Output()
		if err != nil {
			return err
		}
		l.Error(ctx, string(output))
		return fmt.Errorf("unable to proceed with deployment because of uncommitted changes")
	}
	l.Infof(ctx, "\t%s No uncommitted changes!", l.SuccessfulCheck())
	return nil
}

// GetGitSha will return the latest gitSha that will be used to create an application.
func (g *Git) GetGitSha(ctx context.Context, appPath string) (string, error) {
	// Gets latest git sha
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = appPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %v", cmd.String(), string(output))
	}

	return string(output), nil
}

func checkGitVersion(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(out))
	}
	// looks like "git version 2.38.1"
	r := regexp.MustCompile("git version ([0-9.]+)")
	matches := r.FindStringSubmatch(string(out))
	if len(matches) > 0 {
		comparison := semver.Compare("2.28", matches[1])
		return comparison >= 1, nil
	}
	return true, nil
}

func GitMainBranch(branch string) bool {
	switch branch {
	case "main", "master":
		return true
	}

	return false
}

func GetGitBranch(appPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = appPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
