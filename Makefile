SHELL=/bin/bash -o pipefail

.PHONY: build
build:
	go build -o meroxa

PRIVATE_REPOS = github.com/meroxa/meroxa-go
.PHONY: gomod
gomod:
	GOPRIVATE=$(PRIVATE_REPOS) go mod tidy && go mod vendor
