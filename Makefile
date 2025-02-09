# Project variables
BINARY_NAME=mock-ses-api
VERSION?=1.0.0
DOCKER_REGISTRY?=demtech
BUILD_DIR=bin
COVERAGE_DIR=coverage
ENV?=development

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)
GOFILES=$(wildcard *.go)
GOPATH=$(shell go env GOPATH)

# Docker related variables
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_TAG?=latest

# Test variables
TEST_FLAGS=-race -v
COVERAGE_FLAGS=-coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

# Colors for terminal output
YELLOW=\033[0;33m
RED=\033[0;31m
GREEN=\033[0;32m
NC=\033[0m # No Color

.PHONY: all build clean run test test-unit test-e2e test-coverage lint vet format help docker docker-run setup deps


## Build:
build: clean ## Build the project
	@echo "${YELLOW}Building ${BINARY_NAME}...${NC}"
	go build -ldflags="-X main.Version=${VERSION}" -o $(GOBIN)/$(BINARY_NAME) ./cmd/main.go

clean: ## Remove build artifacts
	@echo "${YELLOW}Cleaning...${NC}"
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)

run: build ## Build and run the project
	@echo "${GREEN}Starting ${BINARY_NAME}...${NC}"
	$(GOBIN)/$(BINARY_NAME)

## Development:
deps: ## Download dependencies
	@echo "${YELLOW}Downloading dependencies...${NC}"
	go mod download
	go mod tidy

## Test:
test-setup: ## Prepare for testing
	@echo "${YELLOW}Setting up test environment...${NC}"
	mkdir -p $(COVERAGE_DIR)

unit-test: test-setup ## Run unit tests
	@echo "${YELLOW}Running unit tests...${NC}"
	go test $(TEST_FLAGS) ./test/... ./pkg/...

e2e-test: test-setup ## Run end-to-end tests
	@echo "${YELLOW}Running e2e tests...${NC}"
	go test $(TEST_FLAGS) ./test/...


test-coverage: test-setup ## Run tests with coverage
	@echo "${YELLOW}Running tests with coverage...${NC}"
	go test $(COVERAGE_FLAGS) $(TEST_FLAGS) ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "${GREEN}Coverage report generated at $(COVERAGE_DIR)/coverage.html${NC}"

test: unit-test e2e-test test-coverage ## Run all tests

## Code Quality:
lint: ## Run linter
	@echo "${YELLOW}Running linter...${NC}"
	golangci-lint run

vet: ## Run go vet
	@echo "${YELLOW}Running go vet...${NC}"
	go vet ./...

format: ## Format code
	@echo "${YELLOW}Formatting code...${NC}"
	gofmt -s -w .
	go mod tidy

check: format lint vet test ## Run all code quality checks


## Docker:
docker-build: ## Build docker image
	@echo "${YELLOW}Building docker image...${NC}"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: ## Push docker image
	@echo "${YELLOW}Pushing docker image...${NC}"
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run: ## Run docker container
	@echo "${GREEN}Running docker container...${NC}"
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)


## Help:
help: ## Show this help message
	@echo "Usage: make ${YELLOW}<target>${NC}"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-15s${NC} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help