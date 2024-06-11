# Makefile

protoc:
	cd proto && protoc --go_out=../protogen/golang --go_opt=paths=source_relative \
		--go-grpc_out=../protogen/golang --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=../protogen/golang --grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		./**/*.proto

server:
	go run cmd/server/main.go

gateway:
	go run cmd/client/main.go

run:
	make server && make gateway

down:
	docker-compose down

build:
	docker-compose build

up:
	make down
	docker-compose up -d

up-fresh:
	make build
	make up