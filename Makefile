GO111MODULE=on

all: test

test:
	@go test -mod=vendor -v -race ./...

fmt:
	@go fmt ./...

.PHONY: test fmt