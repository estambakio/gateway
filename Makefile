.DEFAULT_GOAL := help

.PHONY: test
test: ## Run tests
	go test -v -cover ./...

.PHONY: build
build: test ## Build binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/proxy main.go

.PHONY: build-image
build-image: ## Build docker image
	docker build -t estambakio/gateway:master .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' && echo "NOTE: You can find Makefile goals implementation stored in \"./build\" directory"
