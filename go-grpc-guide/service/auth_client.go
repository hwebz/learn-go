package service

import (
	"context"
	pb "github.com/hwebz/go-grpc-guide/pb"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	server   pb.AuthServiceClient
	username string
	password string
}

func NewAuthClient(conn *grpc.ClientConn, username string, password string) *AuthClient {
	service := pb.NewAuthServiceClient(conn)
	return &AuthClient{
		service,
		username,
		password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.server.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
