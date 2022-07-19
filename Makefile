SERVER_NAME="server"
CLIENT_NAME="client"

# .PHONY: build
.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	@ golangci-lint run ./...

test: ## Run unit  tests
	# @go test -short ${PKG_LIST}
	@go test -v -race -timeout 30s ./...

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

build: ## Build the binary file for Intel and ARM architecture
	GOARCH=arm64 GOOS=darwin go build -o ${SERVER_NAME}-darwin -v ./cmd/server
	GOARCH=amd64 GOOS=linux go build -o ${SERVER_NAME} -v ./cmd/server

	GOARCH=arm64 GOOS=darwin go build -o ${CLIENT_NAME}-darwin -v ./cmd/agent
	GOARCH=amd64 GOOS=linux go build -o ${CLIENT_NAME} -v ./cmd/agent

clean: ## Remove previous build
	@rm -f $(SERVER_NAME)
	@rm -f ${SERVER_NAME}-darwin

	@rm -f $(CLIENT_NAME)
	@rm -f ${CLIENT_NAME}-darwin

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build