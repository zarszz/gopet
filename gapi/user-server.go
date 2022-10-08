package gapi

import (
	"go-grpc/config"
	"go-grpc/pb"
	"go-grpc/services"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	config      config.Config
	userService services.UserService
}

func NewGrpcUserServer(config config.Config, userService services.UserService) (*UserServer, error) {
	server := &UserServer{
		config:      config,
		userService: userService,
	}
	return server, nil
}
