.PHONY: gomod
gomod:
	go mod tidy && go mod vendor

.PHONY: test
test:
	go test ./... -coverprofile=c.out -covermode=atomic -v


.PHONY: proto
proto: turbine_proto process_ruby_proto turbine_ruby_proto

.PHONY: turbine_proto
turbine_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto:/defs \
		-v $(CURDIR)/lib/go:/out \
		namely/protoc-all  \
		-f ./turbine/v1/turbine.proto -l go -o /out

.PHONY: process_ruby_proto
process_ruby_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto:/defs \
		-v $(CURDIR)/lib/ruby/turbine_rb/lib/proto:/out \
		namely/protoc-all  \
		-f ./process/v1/service.proto -l ruby -o /out

.PHONY: turbine_ruby_proto
turbine_ruby_proto:
	docker run \
		--rm \
		-v $(CURDIR)/proto:/defs \
		-v $(CURDIR)/lib/ruby/turbine_rb/lib/proto:/out \
		namely/protoc-all  \
		-f ./turbine/v1/turbine.proto -l ruby -o /out
