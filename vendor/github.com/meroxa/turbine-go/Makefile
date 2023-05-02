.PHONY: build install proto test lint gomod

SHELL                = /bin/bash -o pipefail
GO_TEST_FLAGS        = -timeout 5m
GO_TEST_EXTRA_FLAGS ?=
MOCKGEN_VERSION     ?= v1.6.0

build:
	go build -mod=vendor .

install:
	go get -d ./...

proto:
	docker run \
		--rm \
		--platform linux/amd64 \
		-v $(CURDIR)/proto:/defs \
		namely/protoc-all \
		--go-source-relative \
		-f ./service.proto \
		-l go \
		--lint \
		-o .

test:
	go test `go list ./... | grep -v 'turbine-go\/init'` \
		$(GO_TEST_FLAGS) $(GO_TEST_EXTRA_FLAGS) \
		./...

gomod:
	go mod tidy && go mod vendor

lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 5m -v

.PHONY: core_proto
core_proto:
	docker run \
		--rm \
		-v $(CURDIR)/../turbine-core/proto/turbine/v1:/defs \
		-v $(CURDIR)/pkg/proto/core:/out \
		namely/protoc-all  \
		--go-source-relative \
		-f ./turbine.proto \
		-l go \
		-o /out
.PHONY: fmt
fmt: gofumpt
	gofumpt -l -w .

.PHONY: generate
generate: mockgen-install
	go generate ./...

mockgen-install:
	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)

gofumpt:
	go install mvdan.cc/gofumpt@latest
