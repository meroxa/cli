package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
)

type FixtureResource struct {
	File       string
	Collection string
}

type fixtureRecord struct {
	Key       interface{}
	Value     map[string]interface{}
	Timestamp string
}

func (f *FixtureResource) ReadAll(ctx context.Context) ([]*pb.Record, error) {
	b, err := os.ReadFile(f.File)
	if err != nil {
		return nil, err
	}

	var rr []*pb.Record
	if f.Collection == "" {
		// hacky hack because of https://github.com/golang/go/issues/22518
		first := strings.Index(string(b), "[")
		last := strings.LastIndex(string(b), "]")

		substr := b[first : last+1]

		var records []fixtureRecord
		if err := json.Unmarshal(substr, &records); err != nil {
			return nil, err
		}
		for _, r := range records {
			rr = append(rr, wrapRecord(r))
		}
	} else {
		var records map[string][]fixtureRecord
		if err := json.Unmarshal(b, &records); err != nil {
			return nil, err
		}
		for _, r := range records[f.Collection] {
			rr = append(rr, wrapRecord(r))
		}
	}

	return rr, nil
}

func wrapRecord(m fixtureRecord) *pb.Record {
	b, _ := json.Marshal(m.Value)

	ts := timestamppb.New(time.Now())
	if m.Timestamp != "" {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		ts = timestamppb.New(t)
	}

	return &pb.Record{
		Key:       fmt.Sprintf("%v", m.Key),
		Value:     b,
		Timestamp: ts,
	}
}

func PrintRecords(name, collection string, rr []*pb.Record) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Destination %s/%s\n", name, collection)
	fmt.Fprintf(w, "----------------------\n")
	fmt.Fprintln(w, "index\trecord")
	fmt.Fprintln(w, "----\t----")
	for i, r := range rr {
		fmt.Fprintf(w, "%d\t%s\n", i, string(r.Value))
		fmt.Fprintln(w, "----\t----")
	}
	fmt.Fprintf(w, "records written\t%d\n", len(rr))
	w.Flush()
}
