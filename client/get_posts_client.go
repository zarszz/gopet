package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go-grpc/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

type ListPostsClient struct {
	service pb.PostServiceClient
}

func NewListPostsClient(conn *grpc.ClientConn) *ListPostsClient {
	service := pb.NewPostServiceClient(conn)
	return &ListPostsClient{service: service}
}

func (listPostsClient *ListPostsClient) ListPosts(args *pb.GetPostsRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	stream, err := listPostsClient.service.GetPosts(ctx, args)
	if err != nil {
		log.Fatalf("[ListPosts] error : %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("[ListPosts] error when receiving post : %v", err)
		}

		// fmt.Println(res)
		data, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			fmt.Printf("[ListPosts] error when parsing string to json : %v", err)
		}
		fmt.Println(string(data), ",")
	}

}
