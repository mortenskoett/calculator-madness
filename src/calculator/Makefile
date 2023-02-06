PROTOS_DIR := ../../protos
GO_MODULE := "api/bla/pb"

.DEFAULT_GOAL := help

build-protos: ## Build protobuf Go stubs
	protoc --proto_path=$(PROTOS_DIR) --go_out=api --go-grpc_out=api $(PROTOS_DIR)/*.proto

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
