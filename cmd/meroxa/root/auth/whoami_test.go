package auth

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/meroxa/cli/config"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/meroxa-go/pkg/mock"
)

func TestWhoAmIExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	cfg := config.NewInMemoryConfig()

	w := WhoAmI{
		logger: logger,
		client: client,
		config: cfg,
	}

	u := meroxa.User{
		UUID:       "1234-5678-9012",
		Username:   "gbutler",
		Email:      "gbutler@email.io",
		GivenName:  "Joseph",
		FamilyName: "Marcell",
		Verified:   true,
	}

	client.
		EXPECT().
		GetUser(
			ctx,
		).
		Return(&u, nil)

	err := w.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := u.Email

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotUser meroxa.User
	err = json.Unmarshal([]byte(gotJSONOutput), &gotUser)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotUser, u) {
		t.Fatalf("expected \"%v\", got \"%v\"", u, gotUser)
	}

	if w.config.GetString(global.ActorEnv) != u.Email {
		t.Fatalf("expected %q key to be %q", global.ActorEnv, u.Email)
	}

	if w.config.GetString(global.ActorUUIDEnv) != u.UUID {
		t.Fatalf("expected %q key to be %q", global.ActorUUIDEnv, u.UUID)
	}
}
