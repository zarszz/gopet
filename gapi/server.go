package gapi

import (
	"go-grpc/config"
	"go-grpc/pb"
	"go-grpc/services"
	"html/template"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	config         config.Config
	authService    services.AuthService
	userService    services.UserService
	userCollection *mongo.Collection
	template       *template.Template
}

func NewGrpcServer(
	config config.Config, authService services.AuthService, userService services.UserService,
	userCollection *mongo.Collection, template *template.Template) (*Server, error) {
	server := &Server{
		config:         config,
		authService:    authService,
		userService:    userService,
		userCollection: userCollection,
		template:       template,
	}
	return server, nil
}
