SHARED_QUEUE := "../shared/queue"
.DEFAULT_GOAL := help

build: ## Build everything
	make build -C $(SHARED_QUEUE)
	docker build -t viewer-cli -f Dockerfile.cli ../

run: ## Run the viewer cli
	docker run viewer-cli

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
