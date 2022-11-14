package local

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	turbinecore "github.com/meroxa/turbine-core"
	pb "github.com/meroxa/turbine-core/lib/go/github.com/meroxa/turbine/core"
	"github.com/meroxa/turbine-core/platform"
	platform2 "github.com/meroxa/turbine-core/servers/platform"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedTurbineServiceServer
	client      *platform.Client
	processes   map[string]*pb.Process
	resources   []pb.Resource
	deploy      bool
	imageName   string
	config      turbinecore.AppConfig
	appRootPath string
	secrets     map[string]string
	gitSha      string
	appUUID     string
}

func (s *server) Init(ctx context.Context, request *pb.InitRequest) (*emptypb.Empty, error) {
	// load app config (app.json)
	ac, err := turbinecore.ReadAppConfig(request.AppName, request.ConfigFilePath)
	if err != nil {
		log.Printf("error initializeing app; err: %s", err)
		return platform2.Empty(), err
	}
	s.config = ac
	s.appRootPath = request.ConfigFilePath
	return platform2.Empty(), nil
}

func (s *server) GetResource(ctx context.Context, id *pb.NameOrUUID) (*pb.Resource, error) {
	r := pb.Resource{}
	if id.Uuid != "" {
		r.Uuid = id.Uuid
	}
	if id.Name != "" {
		r.Name = id.Name
	}

	s.resources = append(s.resources, r)

	return &r, nil
}

func (s *server) ReadCollection(ctx context.Context, request *pb.ReadCollectionRequest) (*pb.Collection, error) {
	if request.Collection == "" {
		return &pb.Collection{}, fmt.Errorf("please provide a collection name to Records()")
	}

	fixtureFile := s.config.Resources[request.Resource.Name]
	resourceFixturesPath := fmt.Sprintf("%s/%s", s.appRootPath, fixtureFile)
	return readFixtures(resourceFixturesPath, request.Collection)
}

func (s *server) WriteCollectionToResource(ctx context.Context, request *pb.WriteCollectionRequest) (*emptypb.Empty, error) {
	if request.Collection.Name == "" {
		return platform2.Empty(), fmt.Errorf("please provide a collection name to Records()")
	}

	prettyPrintRecords(request.Resource.Name, request.Collection.Stream, request.Collection.Records)

	return platform2.Empty(), nil
}

func (s *server) AddProcessToCollection(ctx context.Context, request *pb.ProcessCollectionRequest) (*pb.Collection, error) {
	s.processes[request.Process.Name] = request.Process
	return nil, nil
}

func New() *server {
	return &server{
		processes: make(map[string]*pb.Process),
		resources: []pb.Resource{},
		deploy:    false,
		imageName: "",
		config:    turbinecore.AppConfig{},
		secrets:   nil,
		gitSha:    "",
		appUUID:   "",
	}
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
