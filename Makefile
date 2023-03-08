SHARED_QUEUE := "./src/shared/queue"
CALC_SERV := "./src/calculator"
VIEW_SERV := "./src/viewer"
.DEFAULT_GOAL := help


build: ## Build everything
	make build -C $(SHARED_QUEUE)
	docker compose build

up: ## Run everything in docker compose
	docker compose up

down: ## Take down everything in docker compose
	docker compose down

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
