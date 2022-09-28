package account

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
	"github.com/spf13/viper"
)

func TestSetAccountsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	a := &meroxa.Account{
		Name: "my-app",
		UUID: "531428f7-4e86-4094-8514-d397d49026f7",
	}

	accounts := []*meroxa.Account{a}

	client.
		EXPECT().
		ListAccounts(ctx).
		Return(accounts, nil)

	cfg := viper.New()
	s := &Set{
		client: client,
		config: cfg,
		logger: logger,
		args:   struct{ NameOrUUID string }{a.UUID},
	}

	err := s.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	if want, got := a.UUID, cfg.GetString(global.UserAccountUUID); want != got {
		t.Fatalf("expected configuration:\n%s\ngot:\n%s", want, got)
	}
}
