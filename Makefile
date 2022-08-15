SHELL           = /bin/bash -o pipefail
GIT_COMMIT     := $(shell git rev-parse --short HEAD)
LDFLAGS        := -X main.GitCommit=${GIT_COMMIT}
GIT_UNTRACKED   = $(shell git diff-index --quiet HEAD -- || echo "(updated)")
LDFLAGS        += -X main.GitUntracked=${GIT_UNTRACKED}
GIT_TAG         = $(shell git describe)
LDFLAGS        += -X main.GitLatestTag=${GIT_TAG}
REBUILD_DOCS    ?= true

.PHONY: build
build: docs
	go build -mod=vendor -o meroxa cmd/meroxa/main.go

.PHONY: install
install:
	go build -ldflags "$(LDFLAGS)" -o $$(go env GOPATH)/bin/meroxa cmd/meroxa/main.go

.PHONY: gomod
gomod:
	go mod tidy && go mod vendor

.PHONY: test
test:
	go test -v ${GO_TEST_FLAGS} -count=1 -timeout 5m ./...

.PHONY: docs
docs:
ifeq ($(REBUILD_DOCS), "true")
	rm -rf docs/cmd && mkdir -p docs/cmd/{md,www}
	rm -rf etc && mkdir -p etc/man/man1 && mkdir -p etc/completion
	go run gen-docs/main.go
endif

.PHONY: lint
lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 5m -v
