gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true --openapiv2_out=swagger

clean:
	rm pb/*.go

run:
	go run main.go

test:
	go test -cover -race ./...

server:
	go run cmd/server/main.go -port 8089

server1:
	go run cmd/server/main.go -port 8090

server-tls:
	go run cmd/server/main.go -port 8089 -tls

server1-tls:
	go run cmd/server/main.go -port 8090 -tls

client:
	go run cmd/client/main.go -address 0.0.0.0:8089

client-tls:
	go run cmd/client/main.go -address 0.0.0.0:8089 -tls

client-lb:
	go run cmd/client/main.go -address 0.0.0.0:8080

client-lb-tls:
	go run cmd/client/main.go -address 0.0.0.0:8080 -tls

evans:
	#go install github.com/ktr0731/evans@latest
	evans -r repl -p 8089

cert:
	cd cert; sh gen.sh; cd ..

rest:
	go run cmd/server/main.go -port 8091 -type rest

.PHONY: gen clean server run client evans test cert