default: help

.PHONY: build

version=`git describe --tags`

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' -e 's/:.*#/ #/'

install: ## Install the binary
	go get -d ./...
	go get -u golang.org/x/lint/golint

build: ## Build the application
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/escher-proxy proxy.go

build-all: ## Build the application
	@echo "version: ${version}"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-x86_64 proxy.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-aarch64 proxy.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-windows-x86_64.exe proxy.go
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-windows-x86.exe proxy.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-darwin-x86_64 proxy.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-darwin-aarch64 proxy.go

run: ## Run the application
	go run proxy.go -v

test: ## Run tests
	go test ./...

lint: ## Check lint errors
	golint -set_exit_status=1 -min_confidence=1.1 ./...