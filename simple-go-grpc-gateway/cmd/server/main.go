package main

import (
	"fmt"
	"github.com/hwebz/simple-go-grpc-gateway/internal"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/orders"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/product"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
)

func main() {
	orderItem := orders.Order{
		OrderId:    10,
		CustomerId: 11,
		IsActive:   true,
		OrderDate:  &date.Date{Year: 2021, Month: 1, Day: 1},
		Products: []*product.Product{
			{
				ProductId:   1,
				ProductName: "CocaCola",
				ProductType: product.ProductType_DRINK,
			},
		},
	}

	// gRPC will typically serialize the message in binary format, which is much faster
	// and take less space compared to a text format like JSON
	bytes, err := protojson.Marshal(&orderItem)
	if err != nil {
		log.Fatal("Deserialization error:", err)
	}

	fmt.Println(string(bytes))

	const addr = "0.0.0.0:50051"

	// create a TCP listener on the specified port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server instance
	server := grpc.NewServer()

	// create an order service instance with a reference to the db
	db := internal.NewDB()
	orderService := internal.NewOrderService(db)

	// register the order service with the grpc server
	orders.RegisterOrdersServer(server, &orderService)

	// start listening to requests
	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
