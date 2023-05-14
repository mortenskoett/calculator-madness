SHARED_QUEUE := "../shared/queue"
PROTOS_DIR := ../../protos

.DEFAULT_GOAL := help

build-protos: ## Build protobuf Go stubs
	protoc --proto_path=$(PROTOS_DIR) --go_out=api --go-grpc_out=api $(PROTOS_DIR)/*.proto

build: ## Build everything
	make build -C $(SHARED_QUEUE)
	docker build -t calculator-server -f Dockerfile.server ../

run: ## Run calculator server
	docker run -p 8000:8000 calculator-server

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
