SHELL=/bin/bash -o pipefail

.PHONY: build
build:
	go build -mod=vendor .

.PHONY: install
install:
	go get -d ./...

.PHONY: proto
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

.PHONY: test
test:
	go test `go list ./... | grep -v init` -timeout 5m ./...

.PHONY: gomod
gomod:
	go mod vendor && go mod tidy
