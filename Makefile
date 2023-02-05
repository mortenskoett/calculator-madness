SRC_DIR := .
PROTO_DIR := internal/protos

.DEFAULT_GOAL := help

build-protos-calculator: ## Build protobuf go stubs
	protoc -I=$(SRC_DIR) --go_out=$(SRC_DIR)/$(PROTO_DIR) $(SRC_DIR)/$(PROTO_DIR)/calculator.proto

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
