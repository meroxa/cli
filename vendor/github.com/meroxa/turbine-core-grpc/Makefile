.PHONY: build generate-clients generate-go generate-python-client generate-ruby-client generate-typescript-client

TURBINE_PROTO="proto/turbine/v1/turbine.proto"
SERVICE_PROTO="proto/process/v1/service.proto"

build: generate-proto
	go build -o turbine-core ./cmd/server/...

generate-proto: generate-go generate-typescript-client generate-ruby-client generate-ruby-service-proto generate-python-client

generate-go:
	protoc --proto_path=proto --go-grpc_out=lib/go --go_out=lib/go $(TURBINE_PROTO)

generate-ruby-client:
	grpc_tools_ruby_protoc -I proto --ruby_out=lib/ruby/turbine/lib/proto --grpc_out=lib/ruby/turbine/lib/proto $(TURBINE_PROTO)

generate-ruby-service-proto:
	grpc_tools_ruby_protoc -I proto --ruby_out=lib/ruby/turbine/lib/proto --grpc_out=lib/ruby/turbine/lib/proto $(SERVICE_PROTO)

generate-python-client:
	protoc --proto_path=proto --python_out=lib/python $(TURBINE_PROTO)

generate-typescript-client:
	protoc --proto_path=proto --ts_out=lib/typescript $(TURBINE_PROTO)

install:
	go build -o $$(go env GOPATH)/bin/turbine-core ./cmd/server/...
