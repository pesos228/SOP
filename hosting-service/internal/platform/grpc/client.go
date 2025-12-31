package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(url string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
