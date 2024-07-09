postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -i POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

.PHONY: postgres createdb dropdb migrateup migratedown

