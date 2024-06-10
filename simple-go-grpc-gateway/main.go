package main

import (
	"fmt"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/orders"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/product"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
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
}
