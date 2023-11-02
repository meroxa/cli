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

func TestListSecrets(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()
	// mockTurbineCLI := turbineMock.NewMockCLI(ctrl)

	list := &List{
		client: client,
		logger: logger,
	}
	body := `
	{
        "page": 1,
        "perPage": 30,
        "totalItems": 2,
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
                }  , 
				{
					"id": "ukip276znrvo2bs",
					"name": "test2",
					"data": {
							"ok": "ok"
					},
					"created": "2023-10-31T18:15:01.169Z",
					"updated": "2023-10-31T18:15:01.169Z"
				}       
        ]
	}`

	httpResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Status:     "200 OK",
		StatusCode: 200,
	}

	client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, nil).Return(
		httpResp,
		nil,
	)

	err := list.Execute(ctx)
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
		if gotSecret.Name == "" {
			t.Fatalf("secret name cannot be empty")
		}
		if fmt.Sprintf("%v", gotSecret.Data) == "" {
			t.Fatalf("secret data cannot be empty")
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
