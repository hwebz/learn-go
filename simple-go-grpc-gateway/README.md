## Simple Go GRPC API Gateway
Reference: https://www.koyeb.com/tutorials/build-a-grpc-api-using-go-and-grpc-gateway

### Install Protobuf Compiler
```shell
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```

### Init go project
```shell
go mod init github.com/hwebz/simple-go-grpc-gateway
```

### Download Google date.proto
```shell
curl -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/type/date.proto -o proto/google/api/date.proto
```

### Install protoc-gen-go
```shell
cd ..
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

### Generate proto files
```shell
make protoc
```

### Install some missing dependencies for generated code
```shell
go get google.golang.org/protobuf # Go implementation for protocol buffers
go get google.golang.org/genproto # Contains the generated Go packages for common protocol buffer types
go get google.golang.org/grpc # for grpc generated codes
```

### Run main application
```shell
make run
```

### Service definition compiler
```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## API Gateway
![img.png](images/api-gateway-diagram.png)

### Binaries need to install for API Gateway
```shell
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```

### Download 2 proto files from Google
```shell
curl -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o proto/google/api/annotations.proto
curl -L https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o proto/google/api/http.proto
```

### Run Gateway and make request
```shell
make server
make gateway

curl -d "@data.json" -X POST -i http://localhost:8080/v0/orders
```

### Sample order payload for adding/updating order
```json
{
    "order": {
        "order_id": 2,
        "customer_id": 1,
        "is_active": true,
        "products": [
          {
            "product_id": "1",
            "product_name": "CocaCola",
            "product_type": "DRINK"
          }
        ],
        "order_date": {
            "year": 2024,
            "month": 6,
            "day": 12
        }
    }
}
```