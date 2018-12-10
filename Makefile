all: test

test:
	@go test -v -race ./...

fmt:
	@go fmt ./...

.PHONY: test fmt