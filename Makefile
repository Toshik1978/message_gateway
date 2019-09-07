GOCMD		= go
GOPATH		= $(shell $(GOCMD) env GOPATH)
GOLINT		= golangci-lint
GOCOV		= gocov
DOCKER		= docker

GIT_VERSION	= $(shell git rev-list -1 HEAD)
CURRENT_DIR = $(shell pwd)

.PHONY: all modules prereq mock build lint test test_race image clean
.DEFAULT_GOAL := all

all: test build

modules:
	@$(GOCMD) mod download

prereq:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@$(GOCMD) get github.com/axw/gocov
	@$(GOCMD) get github.com/golang/mock
	@$(GOCMD) install github.com/axw/gocov/gocov
	@$(GOCMD) install github.com/golang/mock/mockgen

mock:
	@$(GOCMD) generate ./...

build:
	@$(GOCMD) build -ldflags "-X main.GitVersion=$(GIT_VERSION)"

lint:
	@$(GOLINT) run ./... -v

test: lint
	@$(GOCOV) test ./... -v | $(GOCOV) report

test_race: lint
	@$(GOCOV) test ./... -race -v | $(GOCOV) report

image:
	@$(DOCKER) build -t message_gateway .

clean:
	@$(GOCMD) clean
	@rm -f ./message_gateway
