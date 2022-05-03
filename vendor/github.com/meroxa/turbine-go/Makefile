SHELL=/bin/bash -o pipefail

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