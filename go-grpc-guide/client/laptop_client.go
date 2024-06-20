package client

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/hwebz/go-grpc-guide/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{service}
}

func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	//laptop := sample.NewLaptop()
	// laptop.Id = "15f1f70b-aa78-44e5-afcc-c656e42f89eb" // Already exists
	// laptop.Id = "invalid-id"
	//laptop.Id = ""
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	// No timeout
	ctx := context.Background()

	res, err := laptopClient.service.CreateLaptop(ctx, req)
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
}

func (laptopClient *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Printf("Searching laptop with filter: %w", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("Cannot search laptop: %w", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatal("Cannot receive response: %w", err)
		}

		laptop := res.GetLaptop()
		log.Printf("Found laptop with id: %s", laptop.GetId())
		log.Print(" + brand: " + laptop.GetBrand())
		log.Print(" + name: " + laptop.GetName())
		log.Print(" + cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print(" + cpu min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Print(" + ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Printf(" + price: ", laptop.GetPriceUsd())
	}
}

func (laptopClient *LaptopClient) UploadImage(laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("Cannot upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("Cannot send image info: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("Cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Cannot receive response: ", err, stream.RecvMsg(nil))
	}

	log.Printf("Image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("Cannot rate laptop: %v", err)
	}

	waitRepsonse := make(chan error)

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("No more response")
				waitRepsonse <- nil
				return
			}
			if err != nil {
				waitRepsonse <- fmt.Errorf("Cannot receive stream response: %v", err)
				return
			}

			log.Print("Received response ", res)
		}
	}()

	for i, laptopId := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopId,
			Score:    scores[i],
		}

		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("Cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("Request is sent", req)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("Cannot close stream: %v", err)
	}

	err = <-waitRepsonse
	return err
}
