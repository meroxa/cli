package platform

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"syscall"
	"time"

	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/proto"

	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type ProtoWrapper struct {
	ProcessMethod func(context.Context, *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error)
}

func (pw ProtoWrapper) Process(ctx context.Context, record *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
	return pw.ProcessMethod(ctx, record)
}

func ServeFunc(f turbine.Function) error {

	convertedFunc := wrapFrameworkFunc(f.Process)

	fn := struct{ ProtoWrapper }{}
	fn.ProcessMethod = convertedFunc

	addr := os.Getenv("MEROXA_FUNCTION_ADDR")
	if addr == "" {
		return fmt.Errorf("Missing MEROXA_FUNCTION_ADDR env var")
	}

	var g run.Group
	g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM))
	{
		gsrv := grpc.NewServer()
		proto.RegisterFunctionServer(gsrv, fn)

		// health check endpoint
		hsrv := health.NewServer()
		hsrv.SetServingStatus("function", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(gsrv, hsrv)

		g.Add(func() error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			return gsrv.Serve(ln)
		}, func(err error) {
			gsrv.GracefulStop()
		})
	}

	return g.Run()
}

func wrapFrameworkFunc(f func([]turbine.Record) []turbine.Record) func(ctx context.Context, record *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
	return func(ctx context.Context, req *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
		rr := f(protoRecordToValveRecord(req))
		return turbineRecordToProto(rr), nil
	}
}

func protoRecordToValveRecord(req *proto.ProcessRecordRequest) []turbine.Record {
	var rr []turbine.Record

	for _, pr := range req.Records {
		log.Printf("Received  %v", pr)
		vr := turbine.Record{
			Key:       pr.GetKey(),
			Payload:   turbine.Payload(pr.GetValue()),
			Timestamp: time.Unix(pr.GetTimestamp(), 0),
		}
		rr = append(rr, vr)
	}

	return rr
}

func turbineRecordToProto(records []turbine.Record) *proto.ProcessRecordResponse {
	var prr []*proto.Record
	for _, vr := range records {
		pr := proto.Record{
			Key:       vr.Key,
			Value:     string(vr.Payload),
			Timestamp: vr.Timestamp.Unix(),
		}
		prr = append(prr, &pr)
	}
	return &proto.ProcessRecordResponse{Records: prr}
}

type LoggerFunc struct{}

func (lf LoggerFunc) Process(ctx context.Context, req *proto.ProcessRecordRequest) (*proto.ProcessRecordResponse, error) {
	log.Printf("# of records: %d", len(req.Records))
	for _, pr := range req.Records {
		log.Printf("Received: \n\tkey: %+v, \n\tvalue: %+v, \n\ttimestamp: %+v", pr.GetKey(), pr.GetValue(), pr.GetTimestamp())
		log.Printf("Record typeOf %+v", reflect.TypeOf(pr))
	}

	prr := &proto.ProcessRecordResponse{Records: req.Records}

	return prr, nil
}
