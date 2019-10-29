GOPATH		= $(shell go env GOPATH)

GIT_VERSION	= $(shell git rev-list -1 HEAD)
CURRENT_DIR = $(shell pwd)

.PHONY: all modules prereq mock build lint test test+race test+ci image clean
.DEFAULT_GOAL := all

all: test build

modules:
	@go mod download

prereq:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@go get -u github.com/axw/gocov/gocov
	@go get -u github.com/golang/mock/mockgen

mock:
	@go generate ./...

build:
	@go build -ldflags "-X main.GitVersion=$(GIT_VERSION)"

lint:
	@golangci-lint run ./... -v

test: lint
	@go test ./...

test+race: lint
	@go test ./... -race

test+ci: lint
	@go test ./... -coverprofile=coverage.txt -covermode=atomic

image:
	@docker build -t message_gateway .

clean:
	@go clean
	@rm -f ./message_gateway
