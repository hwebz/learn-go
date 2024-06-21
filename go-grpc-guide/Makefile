gen:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative

clean:
	rm pb/*.go

run:
	go run main.go

test:
	go test -cover -race ./...

server:
	go run cmd/server/main.go -port 8089

client:
	go run cmd/client/main.go -address 0.0.0.0:8089

evans:
	#go install github.com/ktr0731/evans@latest
	evans -r repl -p 8089

.PHONY: gen clean server run client evans test