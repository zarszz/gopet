package gapi

import (
	"context"
	"go-grpc/pb"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *PostServer) GetPost(ctx context.Context, request *pb.PostRequest) (*pb.PostResponse, error) {

	post, err := server.postService.FindPostById(request.GetId())
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, err
	}

	postResponse := &pb.PostResponse{
		Post: &pb.Post{
			Id:        post.Id.Hex(),
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			User:      post.User,
			CreatedAt: timestamppb.New(post.CreateAt),
			UpdatedAt: timestamppb.New(post.UpdateAt),
		},
	}

	return postResponse, nil
}
