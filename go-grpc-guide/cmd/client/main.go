package main

import (
	"context"
	"flag"
	pb "github.com/hwebz/go-grpc-guide/pb"
	"github.com/hwebz/go-grpc-guide/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	// laptop.Id = "15f1f70b-aa78-44e5-afcc-c656e42f89eb" // Already exists
	// laptop.Id = "invalid-id"
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("Laptop already exists")
		} else {
			log.Fatal("Cannot create laptop: ", err)
		}
		return
	}

	log.Printf("Created laptop with id: %s", res.Id)
	time.Sleep(time.Second)
}
