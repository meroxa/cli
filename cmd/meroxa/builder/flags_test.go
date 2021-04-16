package builder_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
)

type testFlags struct {
	Flag1  string        `long:"flag1"  short:"a" usage:"flag1 usage"  required:"true" persistent:"false"`
	Flag2  int           `long:"flag2"  short:"b" usage:"flag2 usage"  required:"false" persistent:"true"`
	Flag3  int8          `long:"flag3"  short:"c" usage:"flag3 usage"  required:"true" persistent:"false"`
	Flag4  int16         `long:"flag4"  short:"d" usage:"flag4 usage"  required:"false" persistent:"true"`
	Flag5  int32         `long:"flag5"  short:"e" usage:"flag5 usage"  required:"true" persistent:"false"`
	Flag6  int64         `long:"flag6"  short:"f" usage:"flag6 usage"  required:"false" persistent:"true"`
	Flag7  float32       `long:"flag7"  short:"g" usage:"flag7 usage"  required:"true" persistent:"false"`
	Flag8  float64       `long:"flag8"  short:"h" usage:"flag8 usage"  required:"false" persistent:"true"`
	Flag9  bool          `long:"flag9"  short:"i" usage:"flag9 usage"  required:"true" persistent:"false"`
	Flag10 time.Duration `long:"flag10" short:"j" usage:"flag10 usage" required:"false" persistent:"true"`
	Flag11 []bool        `long:"flag11" short:"k" usage:"flag11 usage" required:"true" persistent:"false"`
	Flag12 []float32     `long:"flag12" short:"l" usage:"flag12 usage" required:"false" persistent:"true"`
	Flag13 []float64     `long:"flag13" short:"m" usage:"flag13 usage" required:"true" persistent:"false"`
	Flag14 []int32       `long:"flag14" short:"n" usage:"flag14 usage" required:"false" persistent:"true"`
	Flag15 []int64       `long:"flag15" short:"o" usage:"flag15 usage" required:"true" persistent:"false"`
	Flag16 []int         `long:"flag16" short:"p" usage:"flag16 usage" required:"false" persistent:"true"`
	Flag17 []string      `long:"flag17" short:"q" usage:"flag17 usage" required:"true" persistent:"false"`
}

func TestBuildFlags(t *testing.T) {
	flags := testFlags{}

	want := []builder.Flag{
		{Long: "flag1", Short: "a", Usage: "flag1 usage", Required: true, Persistent: false, Ptr: &flags.Flag1},
		{Long: "flag2", Short: "b", Usage: "flag2 usage", Required: false, Persistent: true, Ptr: &flags.Flag2},
		{Long: "flag3", Short: "c", Usage: "flag3 usage", Required: true, Persistent: false, Ptr: &flags.Flag3},
		{Long: "flag4", Short: "d", Usage: "flag4 usage", Required: false, Persistent: true, Ptr: &flags.Flag4},
		{Long: "flag5", Short: "e", Usage: "flag5 usage", Required: true, Persistent: false, Ptr: &flags.Flag5},
		{Long: "flag6", Short: "f", Usage: "flag6 usage", Required: false, Persistent: true, Ptr: &flags.Flag6},
		{Long: "flag7", Short: "g", Usage: "flag7 usage", Required: true, Persistent: false, Ptr: &flags.Flag7},
		{Long: "flag8", Short: "h", Usage: "flag8 usage", Required: false, Persistent: true, Ptr: &flags.Flag8},
		{Long: "flag9", Short: "i", Usage: "flag9 usage", Required: true, Persistent: false, Ptr: &flags.Flag9},
		{Long: "flag10", Short: "j", Usage: "flag10 usage", Required: false, Persistent: true, Ptr: &flags.Flag10},
		{Long: "flag11", Short: "k", Usage: "flag11 usage", Required: true, Persistent: false, Ptr: &flags.Flag11},
		{Long: "flag12", Short: "l", Usage: "flag12 usage", Required: false, Persistent: true, Ptr: &flags.Flag12},
		{Long: "flag13", Short: "m", Usage: "flag13 usage", Required: true, Persistent: false, Ptr: &flags.Flag13},
		{Long: "flag14", Short: "n", Usage: "flag14 usage", Required: false, Persistent: true, Ptr: &flags.Flag14},
		{Long: "flag15", Short: "o", Usage: "flag15 usage", Required: true, Persistent: false, Ptr: &flags.Flag15},
		{Long: "flag16", Short: "p", Usage: "flag16 usage", Required: false, Persistent: true, Ptr: &flags.Flag16},
		{Long: "flag17", Short: "q", Usage: "flag17 usage", Required: true, Persistent: false, Ptr: &flags.Flag17},
	}

	got := builder.BuildFlags(&flags)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`expected "%v", got "%v"`, want, got)
	}
}
