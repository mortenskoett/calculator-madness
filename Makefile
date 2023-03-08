SHARED_QUEUE := "./src/shared/queue"
CALC_SERV := "./src/calculator"
VIEW_SERV := "./src/viewer"
.DEFAULT_GOAL := help


build: ## Build everything
	make build -C $(SHARED_QUEUE)
	make build -C $(CALC_SERV)
	make build -C $(VIEW_SERV)

run: ## Run everyting in docker compose
	docker compose down
	docker compose up

help:
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
