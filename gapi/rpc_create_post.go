package gapi

import (
	"context"
	"go-grpc/models"
	"go-grpc/pb"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	post := &models.CreatePostRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
		Image:   req.GetImage(),
		User:    req.GetUser(),
	}

	newPost, err := server.postService.CreatePost(post)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}

	res := &pb.PostResponse{
		Post: &pb.Post{
			Id:        newPost.Id.Hex(),
			Title:     newPost.Title,
			Content:   newPost.Content,
			Image:     newPost.Image,
			User:      newPost.User,
			CreatedAt: timestamppb.New(newPost.CreateAt),
			UpdatedAt: timestamppb.New(newPost.UpdateAt),
		},
	}

	return res, nil
}
