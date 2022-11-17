package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	turbinecore "github.com/meroxa/turbine-core"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type runService struct {
	pb.UnimplementedTurbineServiceServer
	config      turbinecore.AppConfig
	appRootPath string
}

func NewRunService() *runService {
	return &runService{
		config: turbinecore.AppConfig{},
	}
}

func (s *runService) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	// load app config (app.json)
	ac, err := turbinecore.ReadAppConfig(request.AppName, request.ConfigFilePath)
	if err != nil {
		log.Printf("error initializeing app; err: %s", err)
		return empty(), err
	}
	s.config = ac
	s.appRootPath = request.ConfigFilePath
	return empty(), nil
}

func (s *runService) GetResource(ctx context.Context, id *pb.GetResourceRequest) (*pb.Resource, error) {
	return &pb.Resource{
		Name: id.Name,
	}, nil
}

func (s *runService) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.Collection == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to Records()")
	}

	fixtureFile := s.config.Resources[request.Resource.Name]
	resourceFixturesPath := fmt.Sprintf("%s/%s", s.appRootPath, fixtureFile)
	return readFixtures(resourceFixturesPath, request.Collection)
}

func (s *runService) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if request.Collection.Name == "" {
		return empty(), fmt.Errorf("please provide a collection name to Records()")
	}

	prettyPrintRecords(request.Resource.Name, request.Collection.Stream, request.Collection.Records)

	return empty(), nil
}

func (s *runService) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	return request.GetCollection(), nil
}

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func readFixtures(path, collection string) (*pb.Collection, error) {
	log.Printf("fixtures path: %s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return &pb.Collection{}, err
	}

	var records map[string][]fixtureRecord
	err = json.Unmarshal(b, &records)
	if err != nil {
		return &pb.Collection{}, err
	}

	var rr []*pb.Record
	for _, r := range records[collection] {
		rr = append(rr, wrapRecord(r))
	}

	col := &pb.Collection{
		Name:    collection,
		Records: rr,
	}

	return col, nil
}

func wrapRecord(m fixtureRecord) *pb.Record {
	b, _ := json.Marshal(m.Value)

	var ts *timestamppb.Timestamp
	if m.Timestamp == "" {
		ts = timestamppb.New(time.Now())
	} else {
		t, _ := time.Parse(time.RFC3339, m.Timestamp)
		ts = timestamppb.New(t)
	}

	return &pb.Record{
		Key:       m.Key,
		Value:     b,
		Timestamp: ts,
	}
}

func prettyPrintRecords(name string, collection string, rr []*pb.Record) {
	fmt.Printf("=====================to %s (%s) resource=====================\n", name, collection)
	for _, r := range rr {
		payloadVal := string(r.Value)
		fmt.Println(payloadVal)
	}
	fmt.Printf("%d record(s) written\n", len(rr))
}
