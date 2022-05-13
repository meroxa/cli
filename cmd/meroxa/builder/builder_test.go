package builder_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/cobra"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
)

type testCmd struct {
	flagLongFoo string
}

var (
	_ builder.CommandWithDocs        = (*testCmd)(nil)
	_ builder.CommandWithAliases     = (*testCmd)(nil)
	_ builder.CommandWithFlags       = (*testCmd)(nil)
	_ builder.CommandWithSubCommands = (*testCmd)(nil)
)

func (c *testCmd) Usage() string {
	return "cmd1"
}
func (c *testCmd) Docs() builder.Docs {
	return builder.Docs{
		Short:   "short-foo",
		Long:    "long-bar",
		Example: "example-baz",
	}
}
func (c *testCmd) Aliases() []string {
	return []string{"foo", "bar"}
}
func (c *testCmd) Flags() []builder.Flag {
	return []builder.Flag{
		{Long: "long-foo", Short: "l", Usage: "test flag", Required: false, Persistent: false, Ptr: &c.flagLongFoo},
	}
}
func (c *testCmd) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&subCmd{}),
	}
}

type subCmd struct{}

func (c *subCmd) Usage() string {
	return "subCmd"
}

func TestBuildCobraCommand_Structural(t *testing.T) {
	cmd := &testCmd{}

	want := &cobra.Command{
		Use:     "cmd1",
		Aliases: []string{"foo", "bar"},
		Short:   "short-foo",
		Long:    "long-bar",
		Example: "example-baz",
	}
	want.Flags().StringVarP(&cmd.flagLongFoo, "long-foo", "l", "", "test flag")
	want.AddCommand(&cobra.Command{Use: "subCmd"})

	got := builder.BuildCobraCommand(cmd)

	// Since we can't compare functions, we ignore RunE (coming from `buildCommandEvent`)
	got.RunE = nil

	if v := cmp.Diff(got, want, cmpopts.IgnoreUnexported(cobra.Command{})); v != "" {
		t.Fatalf(v)
	}
}

var (
	_ builder.CommandWithArgs     = (*mockCommand)(nil)
	_ builder.CommandWithLogger   = (*mockCommand)(nil)
	_ builder.CommandWithExecute  = (*mockCommand)(nil)
	_ builder.CommandWithoutEvent = (*mockCommand)(nil)
)

// mockCommand is a mock of a behavioral command (5 interfaces).
type mockCommand struct {
	ctrl     *gomock.Controller
	recorder *mockCommandMockRecorder
}

func (m *mockCommand) Event() bool {
	return false
}

// mockCommandMockRecorder is the mock recorder for mockCommand.
type mockCommandMockRecorder struct {
	mock *mockCommand
}

// newMockCommand creates a new mock instance.
func newMockCommand(ctrl *gomock.Controller) *mockCommand {
	mock := &mockCommand{ctrl: ctrl}
	mock.recorder = &mockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *mockCommand) EXPECT() *mockCommandMockRecorder {
	return m.recorder
}

func (m *mockCommand) Usage() string {
	return "mockCmd"
}

// ParseArgs mocks base method.
func (m *mockCommand) ParseArgs(strings []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseArgs", strings)
	ret0, _ := ret[0].(error)
	return ret0
}

// ParseArgs indicates an expected call of ParseArgs.
func (mr *mockCommandMockRecorder) ParseArgs(strings interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseArgs", reflect.TypeOf((*mockCommand)(nil).ParseArgs), strings)
}

// Logger mocks base method.
func (m *mockCommand) Logger(logger log.Logger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Logger", logger)
}

// Logger indicates an expected call of Logger.
func (mr *mockCommandMockRecorder) Logger(logger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*mockCommand)(nil).Logger), logger)
}

// Execute mocks base method.
func (m *mockCommand) Execute(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *mockCommandMockRecorder) Execute(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*mockCommand)(nil).Execute), ctx)
}

func TestBuildCobraCommand_Behavioral(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	cmd := newMockCommand(ctrl)

	// build without setting up mock, we expect no calls when building the cobra command
	got := builder.BuildCobraCommand(cmd)

	// set up expectations before executing the command
	var i int
	cmd.EXPECT().
		ParseArgs(gomock.Any()).
		DoAndReturn(func([]string) error {
			i++
			if i != 1 {
				t.Fatalf("unexpected function order")
			}
			return nil
		})

	cmd.EXPECT().
		Logger(gomock.Any()).
		DoAndReturn(func(log.Logger) {
			i++
			if i != 2 {
				t.Fatalf("unexpected function order")
			}
		})

	cmd.EXPECT().
		Execute(ctx).
		DoAndReturn(func(context.Context) error {
			i++
			if i != 3 {
				t.Fatalf("unexpected function order")
			}
			return nil
		})

	err := got.ExecuteContext(ctx)
	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}
}

type testCmdWithFlags struct {
	flag1  string
	flag2  int
	flag3  int8
	flag4  int16
	flag5  int32
	flag6  int64
	flag7  float32
	flag8  float64
	flag9  bool
	flag10 time.Duration
	flag11 []bool
	flag12 []float32
	flag13 []float64
	flag14 []int32
	flag15 []int64
	flag16 []int
	flag17 []string
}

var _ builder.CommandWithFlags = (*testCmdWithFlags)(nil)

func (t *testCmdWithFlags) Usage() string {
	return "testCmdWithFlags"
}

func (t *testCmdWithFlags) Flags() []builder.Flag {
	return []builder.Flag{
		{Long: "flag1", Short: "a", Usage: "flag1 usage", Required: true, Persistent: false, Ptr: &t.flag1},
		{Long: "flag2", Short: "b", Usage: "flag2 usage", Required: false, Persistent: true, Ptr: &t.flag2},
		{Long: "flag3", Short: "c", Usage: "flag3 usage", Required: true, Persistent: false, Ptr: &t.flag3},
		{Long: "flag4", Short: "d", Usage: "flag4 usage", Required: false, Persistent: true, Ptr: &t.flag4},
		{Long: "flag5", Short: "e", Usage: "flag5 usage", Required: true, Persistent: false, Ptr: &t.flag5},
		{Long: "flag6", Short: "f", Usage: "flag6 usage", Required: false, Persistent: true, Ptr: &t.flag6},
		{Long: "flag7", Short: "g", Usage: "flag7 usage", Required: true, Persistent: false, Ptr: &t.flag7},
		{Long: "flag8", Short: "h", Usage: "flag8 usage", Required: false, Persistent: true, Ptr: &t.flag8},
		{Long: "flag9", Short: "i", Usage: "flag9 usage", Required: true, Persistent: false, Ptr: &t.flag9},
		{Long: "flag10", Short: "j", Usage: "flag10 usage", Required: false, Persistent: true, Ptr: &t.flag10},
		{Long: "flag11", Short: "k", Usage: "flag11 usage", Required: true, Persistent: false, Ptr: &t.flag11},
		{Long: "flag12", Short: "l", Usage: "flag12 usage", Required: false, Persistent: true, Ptr: &t.flag12},
		{Long: "flag13", Short: "m", Usage: "flag13 usage", Required: true, Persistent: false, Ptr: &t.flag13},
		{Long: "flag14", Short: "n", Usage: "flag14 usage", Required: false, Persistent: true, Ptr: &t.flag14},
		{Long: "flag15", Short: "o", Usage: "flag15 usage", Required: true, Persistent: false, Ptr: &t.flag15},
		{Long: "flag16", Short: "p", Usage: "flag16 usage", Required: false, Persistent: true, Ptr: &t.flag16},
		{Long: "flag17", Short: "q", Usage: "flag17 usage", Required: true, Persistent: false, Ptr: &t.flag17},
	}
}

func TestBuildCommandWithFlags(t *testing.T) {
	cmd := &testCmdWithFlags{}

	want := &cobra.Command{Use: "testCmdWithFlags"}
	want.Flags().StringVarP(&cmd.flag1, "flag1", "a", "", "flag1 usage")
	want.PersistentFlags().IntVarP(&cmd.flag2, "flag2", "b", 0, "flag2 usage")
	want.Flags().Int8VarP(&cmd.flag3, "flag3", "c", 0, "flag3 usage")
	want.PersistentFlags().Int16VarP(&cmd.flag4, "flag4", "d", 0, "flag4 usage")
	want.Flags().Int32VarP(&cmd.flag5, "flag5", "e", 0, "flag5 usage")
	want.PersistentFlags().Int64VarP(&cmd.flag6, "flag6", "f", 0, "flag6 usage")
	want.Flags().Float32VarP(&cmd.flag7, "flag7", "g", 0, "flag7 usage")
	want.PersistentFlags().Float64VarP(&cmd.flag8, "flag8", "h", 0, "flag8 usage")
	want.Flags().BoolVarP(&cmd.flag9, "flag9", "i", false, "flag9 usage")
	want.PersistentFlags().DurationVarP(&cmd.flag10, "flag10", "j", 0, "flag10 usage")
	want.Flags().BoolSliceVarP(&cmd.flag11, "flag11", "k", nil, "flag11 usage")
	want.PersistentFlags().Float32SliceVarP(&cmd.flag12, "flag12", "l", nil, "flag12 usage")
	want.Flags().Float64SliceVarP(&cmd.flag13, "flag13", "m", nil, "flag13 usage")
	want.PersistentFlags().Int32SliceVarP(&cmd.flag14, "flag14", "n", nil, "flag14 usage")
	want.Flags().Int64SliceVarP(&cmd.flag15, "flag15", "o", nil, "flag15 usage")
	want.PersistentFlags().IntSliceVarP(&cmd.flag16, "flag16", "p", nil, "flag16 usage")
	want.Flags().StringSliceVarP(&cmd.flag17, "flag17", "q", nil, "flag17 usage")

	for i := 1; i <= 17; i++ {
		if i%2 == 1 {
			_ = want.MarkFlagRequired(fmt.Sprintf("flag%d", i))
		}
	}

	got := builder.BuildCobraCommand(cmd)

	// Since we can't compare functions, we ignore RunE (coming from `buildCommandEvent`)
	got.RunE = nil

	// Since we can't compare functions, we ignore PostRunE (coming from `buildCommandAutoUpdate`)
	got.PostRunE = nil

	if v := cmp.Diff(got, want, cmpopts.IgnoreUnexported(cobra.Command{})); v != "" {
		t.Fatalf(v)
	}
}
