package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

// This variable is passed from Dockerfile as environment variable
var orderServiceAddr string

func main() {
	// Set up a connection to the order server.
	fmt.Println("Connecting to order service via", orderServiceAddr)
	conn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to order service: %v", err)
	}
	defer conn.Close()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	if err = orders.RegisterOrdersHandler(context.Background(), mux, conn); err != nil {
		log.Fatalf("Failed to register the order server: %v", err)
	}

	// start listening to requests from the gateway server
	addr := "0.0.0.0:8080"
	fmt.Println("API Gateway server is running on " + addr)
	if err = http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Gateway server closed abruptly: ", err)
	}
}
