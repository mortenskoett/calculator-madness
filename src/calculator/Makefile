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

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
