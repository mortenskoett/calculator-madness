SHARED_QUEUE := "../shared/queue"
PROTOS_DIR := ../../protos

.DEFAULT_GOAL := help

build-protos: ## Build protobuf Go stubs
	protoc --proto_path=$(PROTOS_DIR) --go_out=api --go-grpc_out=api $(PROTOS_DIR)/*.proto

build-cli: ## Build cli
	make build -C $(SHARED_QUEUE)
	docker build -t viewer-cli -f Dockerfile.cli ../

run-cli: ## Run viewer cli
	docker run viewer-cli

build-web: ## Build web server
	make build -C $(SHARED_QUEUE)
	docker build -t viewer-web -f Dockerfile.web ../

run-web: ## Run web server
	docker run -p 3000:3000 viewer-web

rerun-web-local: ## Rerun web app on changes
	reflex -s -- sh -c 'go run ./cmd/web'

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
