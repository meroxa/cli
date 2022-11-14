.PHONY: gomod
gomod:
	go mod tidy && go mod vendor

.PHONY: test
test:
	go test ./... -coverprofile=c.out -covermode=atomic -v
