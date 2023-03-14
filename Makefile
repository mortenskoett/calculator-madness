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

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
