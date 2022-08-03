.PHONY: build install proto test lint gomod

SHELL                = /bin/bash -o pipefail
GO_TEST_FLAGS        = -timeout 5m
GO_TEST_EXTRA_FLAGS ?=

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
	go mod vendor && go mod tidy

lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 5m -v
