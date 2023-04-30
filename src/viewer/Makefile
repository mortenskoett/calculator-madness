SHARED_QUEUE := "../shared/queue"
.DEFAULT_GOAL := help

build: ## Build everything
	make build -C $(SHARED_QUEUE)
	docker build -t viewer-cli -f Dockerfile.cli ../

run: ## Run the viewer cli
	docker run viewer-cli

rerun-web: ## Rerun web app on changes
	reflex -s -- sh -c 'go run ./cmd/web'

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
