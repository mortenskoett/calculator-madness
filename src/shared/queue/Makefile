.DEFAULT_GOAL := help

build: ## Build shared lib
	docker build -t calculator-shared-queue -f Dockerfile.queue ../../

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
