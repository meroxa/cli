package internal

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FixtureResource struct {
	File       string
	Collection string
}

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func (f *FixtureResource) ReadAll(ctx context.Context) ([]*pb.Record, error) {
	b, err := os.ReadFile(f.File)
	if err != nil {
		return nil, err
	}

	var records map[string][]fixtureRecord
	if err := json.Unmarshal(b, &records); err != nil {
		return nil, err
	}

	var rr []*pb.Record
	for _, r := range records[f.Collection] {
		rr = append(rr, wrapRecord(r))
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
		Key:       m.Key,
		Value:     b,
		Timestamp: ts,
	}
}

func PrintRecords(name, collection string, rr []*pb.Record) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Destination %s/%s", name, collection)
	t.AppendHeader(table.Row{"index", "record"})
	for i, r := range rr {
		t.AppendRow(table.Row{i, string(r.Value)})
	}
	t.AppendFooter(table.Row{"records written", len(rr)})
	t.Render()
}
