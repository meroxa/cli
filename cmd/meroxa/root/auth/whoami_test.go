package auth

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
)

func TestWhoAmIExecution(t *testing.T) {
	ctx := context.Background()
	logger := log.NewTestLogger()
	email := "user@example.com"
	os.Setenv(global.TenantEmailAddress, email)

	cfg := config.NewInMemoryConfig()

	w := WhoAmI{
		logger: logger,
		config: cfg,
	}

	err := w.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := email

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := strings.TrimSpace(logger.JSONOutput())
	if gotJSONOutput != email {
		t.Fatalf("expected \"%v\", got \"%v\"", email, gotJSONOutput)
	}
}
