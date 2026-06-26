include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d workoutapp-postgres

env-down:
	@dokcer compose down workoutapp-postgres

env-cleanup:
	@read -p "Очистить все данные DB? [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down -v && \
		echo "Окружение очищено (volumes удалены)"; \
	else \
		echo "Отменено"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутсвует необходимый параметр seq"; \
		exit 1; \
	fi; \
	docker compose run --rm workoutapp-postgres-migrate \
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
		"$(action)"
