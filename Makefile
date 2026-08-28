include .env
export

export PROJECT_ROOT=${shell pwd}

docker-up:
	@docker compose up -d postgres-service port-forwarder

docker-down:
	@docker compose down

volume-cleanup:
	@read -p "Clean all volume data? Danger data lost. [Y/n]: " ans; \
	if [ "$$ans" = "Y" ]; then \
		docker compose down postgres-service port-forwarder && \
		docker volume rm todolist-golang_todoapp_postgres_data && \
		echo "Data files cleaned"; \
	else \
		echo "Cleanup canceled"; \
	fi

logs-cleanup:
	@read -p "Clean all logs data? Danger data lost. [Y/n]: " ans; \
	if [ "$$ans" = "Y" ]; then \
		rm -rf ${PROJECT_ROOT}/data/logs/ && \
		echo "Data files cleaned"; \
	else \
		echo "Cleanup canceled"; \
	fi

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Variable 'name' not found. Example: make migrate-create name=init"; \
		exit 1; \
	fi; \

	@docker compose run --rm migrate-service \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(name)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Variable 'action' not found. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \

	@docker compose run --rm migrate-service \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-service:$(POSTGRES_PORT)/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

migrate-up:
	@make migrate-action action=up
	
migrate-down:
	@make migrate-action action=down

port-forward-start:
	@docker compose up -d port-forwarder

port-forward-stop:
	@docker compose down port-forwarder

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/data/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/api/main.go

generate-secret:
	openssl rand -base64 32