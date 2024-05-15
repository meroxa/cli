package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
	"github.com/meroxa/cli/log"
)

func TestCreateSecrets(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()
	// mockConduitCLI := ConduitMock.NewMockCLI(ctrl)

	create := &Create{
		client: client,
		logger: logger,
		flags: struct {
			Data string "long:\"data\" usage:\"Secret's data, passed as a JSON string\""
		}{Data: `{"ok": "ok"}`},
		args: struct {
			secretName string
		}{secretName: "test"},
	}
	body := `{
		"collectionId": "pmm1jxdx100l3ux",
		"collectionName": "secrets",
		"created": "2023-11-01 06:12:12.682Z",
		"data": {
		  "ok": "ok"
		},
		"id": "o7kh2ekz3rdagrv",
		"name": "test",
		"updated": "2023-11-01 06:12:12.682Z"
	  }`

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}

	secret := &Secrets{
		Name: "test",
		Data: map[string]interface{}{
			"ok": "ok",
		},
	}
	client.EXPECT().CollectionRequest(ctx, "POST", collectionName, "", secret, nil).Return(
		httpResp,
		nil,
	)

	err := create.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotJSONOutput := logger.JSONOutput()

	var gotSecret Secrets
	err = json.Unmarshal([]byte(gotJSONOutput), &gotSecret)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if gotSecret.Name != secret.Name {
		t.Fatalf("expected \"%s\" got \"%s\"", gotSecret.Name, secret.Name)
	}
	if fmt.Sprintf("%v", gotSecret.Data) != fmt.Sprintf("%v", secret.Data) {
		t.Fatalf("expected \"%s\" got \"%s\"", gotSecret.Name, secret.Name)
	}
	if gotSecret.ID == "" {
		t.Fatalf("secret ID cannot be empty")
	}
	if gotSecret.Created.String() == "" {
		t.Fatalf("secret created time cannot be empty")
	}
	if gotSecret.Updated.String() == "" {
		t.Fatalf("secret updated time cannot be empty")
	}
}
