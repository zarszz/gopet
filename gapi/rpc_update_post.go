package gapi

import (
	"context"
	"go-grpc/models"
	"go-grpc/pb"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *PostServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error) {
	id := req.GetId()

	post := &models.UpdatePost{
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		Image:    req.GetImage(),
		User:     req.GetUser(),
		UpdateAt: time.Now(),
	}

	updatedPost, err := server.postService.UpdatePost(id, post)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, err
	}

	res := &pb.PostResponse{
		Post: &pb.Post{
			Id:        updatedPost.Id.Hex(),
			Title:     updatedPost.Title,
			Content:   updatedPost.Content,
			Image:     updatedPost.Image,
			User:      updatedPost.User,
			CreatedAt: timestamppb.New(updatedPost.CreateAt),
			UpdatedAt: timestamppb.New(updatedPost.UpdateAt),
		},
	}

	return res, nil
}
