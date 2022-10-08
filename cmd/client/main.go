package main

import (
	"go-grpc/client"
	"go-grpc/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	// get lists post
	listPostClient := client.NewListPostsClient(conn)
	var page int64 = 1
	var limit int64 = 10
	args := &pb.GetPostsRequest{
		Page:  &page,
		Limit: &limit,
	}

	listPostClient.ListPosts(args)
}
