GOCMD		= go
GOPATH		= $(shell $(GOCMD) env GOPATH)
GOLINT		= golangci-lint
GOCOV		= gocov
DOCKER		= docker

GIT_VERSION	= $(shell git rev-list -1 HEAD)
CURRENT_DIR = $(shell pwd)

.PHONY: all prereq build lint test mock clean
.DEFAULT_GOAL := all

all: test build

prereq:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@$(GOCMD) get github.com/axw/gocov/gocov
	@$(GOCMD) get github.com/golang/mock/gomock
	@$(GOCMD) install github.com/golang/mock/mockgen

mock:
	@$(GOCMD) generate ./...

build:
	@$(GOCMD) build -ldflags "-X main.GitVersion=$(GIT_VERSION)"

lint:
	@$(GOLINT) run ./... -E stylecheck -E gofmt -E goimports -E golint -E unconvert -E goconst -E unparam -E scopelint -E lll -v --skip-dirs "mock" --tests=false

test: lint
	@$(GOCOV) test ./... -v | $(GOCOV) report

test_race: lint
	@$(GOCOV) test ./... -race -v | $(GOCOV) report

clean:
	@$(GOCMD) clean
	@rm -f ./message_gateway
