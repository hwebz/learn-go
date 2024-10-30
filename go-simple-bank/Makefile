postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=myuser --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://myuser:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://myuser:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

# Define a variable to handle both cases
DOCKER_COMPOSE = $(shell if docker compose version >/dev/null 2>&1; then echo 'docker compose'; else echo 'docker-compose'; fi)

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

build:
	$(DOCKER_COMPOSE) build

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown up down build sqlc

