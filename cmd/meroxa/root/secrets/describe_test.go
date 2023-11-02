package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	basicMock "github.com/meroxa/cli/cmd/meroxa/global/mock"
	"github.com/meroxa/cli/log"
)

func TestDescribeSecrets(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()
	// mockTurbineCLI := turbineMock.NewMockCLI(ctrl)

	describe := &Describe{
		client: client,
		logger: logger,
		args:   struct{ nameOrUUID string }{nameOrUUID: "test"},
	}
	body := `
	{
        "page": 1,
        "perPage": 30,
        "totalItems": 1,
        "totalPages": 1,
        "items": [
                {
                        "id": "ukip276znrvo2bs",
                        "name": "test",
                        "data": {
                                "ok": "ok"
                        },
                        "created": "2023-10-31T18:15:01.169Z",
                        "updated": "2023-10-31T18:15:01.169Z"
                }           
        ]
}
	`

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

	a := &url.Values{}
	a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", secret.Name, secret.Name))

	client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, *a).Return(
		httpResp,
		nil,
	)

	err := describe.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	gotJSONOutput := logger.JSONOutput()

	var listSecrets ListSecrets
	err = json.Unmarshal([]byte(gotJSONOutput), &listSecrets)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	for _, gotSecret := range listSecrets.Items {
		if gotSecret.Name != secret.Name {
			t.Fatalf("expected \"%s\" got \"%s\"", gotSecret.Name, secret.Name)
		}
		if fmt.Sprintf("%v", gotSecret.Data) != fmt.Sprintf("%v", secret.Data) {
			t.Fatalf("expected \"%s\" got \"%s\"", fmt.Sprintf("%v", gotSecret.Data), fmt.Sprintf("%v", secret.Data))
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
}
