SHELL=/bin/bash -o pipefail

.PHONY: build
build: docs
	go build -mod=vendor -o meroxa

.PHONY: install
install:
	go build -o $$(go env GOPATH)/bin/meroxa

PRIVATE_REPOS = github.com/meroxa/meroxa-go
.PHONY: gomod
gomod:
	GOPRIVATE=$(PRIVATE_REPOS) go mod tidy && go mod vendor

.PHONY: test
test:
	go test -v ${GO_TEST_FLAGS} -count=1 -timeout 5m ./...

.PHONY: docs
docs:
	rm -rf docs/cmd && mkdir -p docs/cmd/{md,www}
	rm -rf etc && mkdir -p etc/man/man1 && mkdir -p etc/completion
	go run gen-docs/main.go