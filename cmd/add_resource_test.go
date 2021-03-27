package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	utils "github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"testing"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err error
		name string
	}{
		{[]string{""},nil, ""},
		{[]string{"resName"},nil, "resName"},
	}

	for _, tt := range tests {
		name, err := AddResource{}.checkArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, name)
		}
	}
}

func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name string
		required bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"credentials", false, ""},
		{"metadata", false, "m"},
	}

	c := &cobra.Command{}
	AddResource{}.setFlags(c)

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
}

func TestAddResourceOutput(t *testing.T) {
	r := utils.GenerateResource()

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	expected := fmt.Sprintf("Resource %s successfully added!", r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestAddResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		AddResource{}.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)


	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}

func TestAddResourceExecution(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	client := NewMockAddResourceClient(ctrl)

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
	}

	client.
		EXPECT().
		CreateResource(
			ctx,
			gomock.Eq(&r),
		).
		DoAndReturn(func() (*meroxa.Resource, error) {
			nr := utils.GenerateResource()
			return &nr, nil
		})

	got, err := AddResource{}.execute(ctx, client, r)

	if got != nil {
		t.Fatal("not good")
	}
	if err == nil {
		t.Fatal("not good")
	}
}

type MockAddResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAddResourceClientMockRecorder
}

// MockAddResourceClientMockRecorder is the mock recorder for MockAddResourceClient.
type MockAddResourceClientMockRecorder struct {
	mock *MockAddResourceClient
}

// NewMockAddResourceClient creates a new mock instance.
func NewMockAddResourceClient(ctrl *gomock.Controller) *MockAddResourceClient {
	mock := &MockAddResourceClient{ctrl: ctrl}
	mock.recorder = &MockAddResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAddResourceClient) EXPECT() *MockAddResourceClientMockRecorder {
	return m.recorder
}

// CreateResource mocks base method.
func (m *MockAddResourceClient) CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateResource", ctx, resource)
	ret0, _ := ret[0].(*meroxa.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateResource indicates an expected call of CreateResource.
func (mr *MockAddResourceClientMockRecorder) CreateResource(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateResource", reflect.TypeOf((*MockAddResourceClient)(nil).CreateResource), ctx, resource)
}

