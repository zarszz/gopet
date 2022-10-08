package gapi

import (
	"go-grpc/config"
	"go-grpc/pb"
	"go-grpc/services"
	"html/template"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	config         config.Config
	authService    services.AuthService
	userService    services.UserService
	userCollection *mongo.Collection
	template       *template.Template
}

func NewGrpcAuthServer(
	config config.Config, authService services.AuthService, userService services.UserService,
	userCollection *mongo.Collection, template *template.Template) (*AuthServer, error) {
	server := &AuthServer{
		config:         config,
		authService:    authService,
		userService:    userService,
		userCollection: userCollection,
		template:       template,
	}
	return server, nil
}
