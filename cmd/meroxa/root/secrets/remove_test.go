package secrets

import (
	"context"
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

func TestRemoveSecrets(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := basicMock.NewMockBasicClient(ctrl)
	logger := log.NewTestLogger()
	// mockTurbineCLI := turbineMock.NewMockCLI(ctrl)

	remove := &Remove{
		client: client,
		logger: logger,
		args:   struct{ nameOrUUID string }{nameOrUUID: "test"},
		flags: struct {
			Force bool "long:\"force\" short:\"f\" default:\"false\" usage:\"skip confirmation\""
		}{Force: true},
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

	a := &url.Values{}
	a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", remove.args.nameOrUUID, remove.args.nameOrUUID))

	client.EXPECT().CollectionRequest(ctx, "GET", collectionName, "", nil, *a).Return(
		httpResp,
		nil,
	)

	client.EXPECT().CollectionRequest(ctx, "DELETE", collectionName, "ukip276znrvo2bs", nil, nil).Return(
		&http.Response{
			Status:     "204 OK",
			StatusCode: 204,
		},
		nil,
	)

	err := remove.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}
}
