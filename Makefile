# Go files to format
BIN = teeworlds-prometheus-exporter
GOFMT_FILES ?= $(shell find . -name "*.go")

default: fmt

fmt:
	gofmt -w $(GOFMT_FILES)

build:
	go build -o $(BIN)

clean:
	go clean -testcache

test: clean
	go test -v ./...

.PHONY: \
	fmt \
	test \
	clean
