.DEFAULT_GOAL := help

build: ## Build shared lib
	docker build -t calculator-shared-queue -f Dockerfile.queue ../../

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
