package gapi

import (
	"context"
	"go-grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userServer *UserServer) GetMe(ctx context.Context, request *pb.GetMeRequest) (*pb.UserResponse, error) {
	id := request.GetId()

	user, err := userServer.userService.FindUserById(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not get user")
	}

	res := &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}

	return res, nil
}
