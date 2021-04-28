package root

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/cmd/meroxa/root/deprecated"

	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	utils "github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
)

func TestWhoAmIExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockGetUserClient(ctrl)

	u := generateUser()

	client.
		EXPECT().
		GetUser(
			ctx,
		).
		Return(&u, nil)

	ar := &GetUser{}
	got, err := ar.execute(ctx, client)

	if !reflect.DeepEqual(got, &u) {
		t.Fatalf("expected \"%v\", got \"%v\"", &u, got)
	}

	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}
}

func TestWhoAmIOutput(t *testing.T) {
	u := generateUser()
	deprecated.FlagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		gu := &GetUser{}
		gu.output(&u)
	})

	expected := u.Email

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestWhoAmIOutputJSONOutput(t *testing.T) {
	u := generateUser()
	deprecated.FlagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		gu := &GetUser{}
		gu.output(&u)
	})

	var parsedOutput meroxa.User
	_ = json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(u, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}

func generateUser() meroxa.User {
	return meroxa.User{
		UUID:       "1234-5678-9012",
		Username:   "gbutler",
		Email:      "gbutler@email.io",
		GivenName:  "Joseph",
		FamilyName: "Marcell",
		Verified:   true,
	}
}
