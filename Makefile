DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

# Define pkgs, run, and cover vairables for test so that we can override them in
# the terminal more easily.
pkgs := $(shell go list ./...)
run := .
count := 1

## help: Show this help message
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help

## clean: clean testcache
clean:
	@echo "--> Clearing testcache"
	@go clean --testcache
.PHONY: clean

## cover: generate to code coverage report.
cover:
	@echo "--> Generating Code Coverage"
	@go install github.com/ory/go-acc@latest
	@go-acc -o coverage.txt $(pkgs)
.PHONY: cover

## deps: Install dependencies
deps:
	@echo "--> Installing dependencies"
	@go mod download
#	@make proto-gen
	@go mod tidy
.PHONY: deps

## lint: Run linters golangci-lint and markdownlint.
lint: vet
	@echo "--> Running golangci-lint"
	@golangci-lint run
	@echo "--> Running markdownlint"
	@markdownlint --config .markdownlint.yaml '**/*.md'
	@echo "--> Running hadolint"
	@hadolint docker/mockserv.Dockerfile
	@echo "--> Running yamllint"
	@yamllint --no-warnings . -c .yamllint.yml

.PHONY: lint

## fmt: Run fixes for linters. Currently only markdownlint.
fmt:
	@echo "--> Formatting markdownlint"
	@markdownlint --config .markdownlint.yaml '**/*.md' -f
.PHONY: fmt

## vet: Run go vet
vet: 
	@echo "--> Running go vet"
	@go vet $(pkgs)
.PHONY: vet

## test: Running unit tests
test: vet
	@echo "--> Running unit tests"
	@go test -v -race -covermode=atomic -coverprofile=coverage.txt $(pkgs) -run $(run) -count=$(count)
.PHONY: test

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

## check-proto-deps: Check protobuf deps
check-proto-deps:
ifeq (,$(shell which protoc-gen-gocosmos))
	@go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
endif
.PHONY: check-proto-deps

## proto-gen: Generate protobuf files
proto-gen: check-proto-deps
	@echo "--> Generating Protobuf files"
	@go run github.com/bufbuild/buf/cmd/buf@latest generate --path proto/execution
.PHONY: proto-gen

## proto-lint: Lint protobuf files.
proto-lint: check-proto-deps
	@echo "--> Linting Protobuf files"
	@go run github.com/bufbuild/buf/cmd/buf@latest lint --error-format=json
.PHONY: proto-lint

## mock-gen: Re-generates DA mock
mock-gen:
	@mockery
.PHONY: mock-gen
