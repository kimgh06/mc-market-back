ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

dev-docker:
	@echo Starting development compose
	docker compose -f docker-compose.development.yml build maple-core --no-cache
	docker compose -f docker-compose.development.yml up -d
	@echo Shutdown development compose

dev-postgres-standalone:
	@echo Starting development PostgreSQL as standalone container
	docker compose -f docker-compose.postgres.yml up -d --remove-orphans
	@echo Started in background

dev-postgres-standalone-stop:
	@echo Stopping standalone development PostgreSQL container
	docker compose -f docker-compose.postgres.yml down
	@echo Stopped

dev-postgres-standalone-reset:
	@echo Recreating PostgreSQL container
	@docker compose -f docker-compose.postgres.yml rm -s -f -v maple-postgres-standalone
	@make dev-postgres-standalone
	@echo Recreated

dev-run:
	go build -o maple
	./maple

dev-generate:
	go build -o maple
	./maple generate

dev-logs:
	@echo Inspecting compose logs
	docker compose -f docker-compose.development.yml --ansi=always logs maple-core

dev-logsnocolor:
	@echo Inspecting compose logs without color
	docker compose -f docker-compose.development.yml logs

migrate-new:
	migrate create -ext sql -dir $(ROOT_DIR)/schema/migrations -seq $(name)
