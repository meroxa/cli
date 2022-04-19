SHELL=/bin/bash -o pipefail

.PHONY: build
build: docs
	go build -mod=vendor -o meroxa cmd/meroxa/main.go

.PHONY: install
install:
	go build -o $$(go env GOPATH)/bin/meroxa cmd/meroxa/main.go

.PHONY: gomod
gomod:
	go mod tidy && go mod vendor

.PHONY: test
test:
	go test -v ${GO_TEST_FLAGS} -count=1 -timeout 5m ./...

.PHONY: docs
docs:
	rm -rf docs/cmd && mkdir -p docs/cmd/{md,www}
	rm -rf etc && mkdir -p etc/man/man1 && mkdir -p etc/completion
	go run gen-docs/main.go

.PHONY: lint
lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 5m -v