GO ?= $(shell command -v go 2> /dev/null)

all: lint generate test run

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	$(GO) test -v -race ./...

.PHONY: run
run:
	$(GO) run ./cmd/generate

.PHONY: build
build:
	$(GO) build -o bin/ -trimpath ./...

#install_mockery:
#	mockery --version | grep v2.32.0 || $(GO) install github.com/vektra/mockery/v2/...@v2.32.0
#
#.PHONY: generate
#generate: install_mockery
#	$(GO) generate ./...

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: install_proto
install_proto:
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0

.PHONY: gen-proto
gen-proto: install_proto
	protoc -I=./protos --go_out=. --go_opt=paths=source_relative \
		protos/*.proto
