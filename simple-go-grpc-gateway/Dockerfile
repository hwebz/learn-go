# Build go services + gateway
FROM golang:1.22-alpine AS builder
ARG ORDER_SERVICE_ADDRESS
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./orders-service ./cmd/server/main.go
RUN go build -ldflags "-X main.orderServiceAddr=$ORDER_SERVICE_ADDRESS" -o ./gateway-service ./cmd/client/main.go

# Run Order service on port 50051
FROM alpine:latest AS orders-service
WORKDIR /app
COPY --from=builder /app/orders-service .
EXPOSE 50051
ENTRYPOINT ["./orders-service"]

# Run gateway on port 8080
FROM alpine:latest AS gateway-service
WORKDIR /app
COPY --from=builder /app/gateway-service .
EXPOSE 8080
ENTRYPOINT ["./gateway-service"]