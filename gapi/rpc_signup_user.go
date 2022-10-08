package gapi

import (
	"context"
	"go-grpc/models"
	"go-grpc/pb"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.GenericResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.SignUpInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	newUser, err := server.authService.SignUpUser(&user, server.template)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	message := "We sent an email with a verification code to " + newUser.Email

	res := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return res, nil
}
