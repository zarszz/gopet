package gapi

import (
	"context"
	"go-grpc/pb"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *PostServer) DeletePost(ctx context.Context, request *pb.PostRequest) (*pb.DeletePostResponse, error) {
	err := server.postService.DeletePost(request.GetId())
	if err != nil {
		if strings.Contains(err.Error(), "exist") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, err
	}

	response := &pb.DeletePostResponse{
		Success: true,
	}

	return response, nil
}
