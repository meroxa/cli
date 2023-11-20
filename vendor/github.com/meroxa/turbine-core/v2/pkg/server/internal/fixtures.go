package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"github.com/meroxa/turbine-core/v2/proto/turbine/v2"
)

func ReadFixture(ctx context.Context, file string) ([]*opencdcv1.Record, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var fixtureRecords []opencdc.Record

	if err := json.Unmarshal(b, &fixtureRecords); err != nil {
		return nil, err
	}

	protoRecords := make([]*opencdcv1.Record, len(fixtureRecords))

	for i, r := range fixtureRecords {
		protoRecords[i] = &opencdcv1.Record{}
		if err := r.ToProto(protoRecords[i]); err != nil {
			return nil, err
		}
	}

	return protoRecords, nil
}

func PrintRecords(name string, sr *turbinev2.StreamRecords) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Destination %s\n", name)
	fmt.Fprintf(w, "----------------------\n")
	fmt.Fprintln(w, "index\trecord")
	fmt.Fprintln(w, "----\t----")

	for i, proto := range sr.Records {
		var r opencdc.Record

		if err := r.FromProto(proto); err != nil {
			fmt.Fprintf(w, "%d\t%s\n", i, fmt.Sprintf("failed to render, error: %s", err.Error()))
		} else {
			fmt.Fprintf(w, "%d\t%s\n", i, string(r.Bytes()))
		}

		fmt.Fprintln(w, "----\t----")
	}

	fmt.Fprintf(w, "records written\t%d\n", len(sr.Records))
	w.Flush()
}
