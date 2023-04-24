SHARED_QUEUE := "./src/shared/queue"
CALC_SERV := "./src/calculator"
VIEW_SERV := "./src/viewer"
.DEFAULT_GOAL := help

run: build ## Build and run everything in docker compose
	docker compose up

build: ## Build everything
	docker compose down
	make build -C $(SHARED_QUEUE) && \
	docker compose build

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
