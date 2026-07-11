include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@mkdir -p out/pgdata
	@docker compose up -d workoutapp-postgres

env-down:
	@docker compose down workoutapp-postgres

env-cleanup:
	@read -p "Очистить все данные DB? [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down workoutapp-postgres && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Отменено"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq"; \
		exit 1; \
	fi
	@mkdir -p migrations
	@docker compose run --rm workoutapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутсвует необходимый параметр action"; \
		exit 1; \
	fi; \
	docker compose run --rm workoutapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@workoutapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action) $(steps)

app-run:
	@export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/main.go
