package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestGetLatestCLITag(t *testing.T) {
	version := "2.2.0"

	f := fmt.Sprintf(`
class Meroxa < Formula
  desc "The Meroxa CLI"
  homepage "https://meroxa.io"
  version "%s"
end
`, version)

	sEnc := base64.StdEncoding.EncodeToString([]byte(f))
	jsonResponse := fmt.Sprintf(`
	{
	   "content": "%s"
	}
	`, sEnc)

	r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	ctx := context.Background()
	gotContent, err := getContentHomebrewFormula(ctx)
	if err != nil {
		t.Error("unexpected error, got: ", err)
		return
	}

	if !reflect.DeepEqual(gotContent, f) {
		t.Errorf("expected %v, got %v", f, gotContent)
	}

	gotVersion := parseVersionFromFormulaFile(gotContent)
	if gotVersion != version {
		t.Errorf("expected version %q, got %q", version, gotVersion)
	}
}
