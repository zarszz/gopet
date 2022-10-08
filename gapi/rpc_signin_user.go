package gapi

import (
	"context"
	"go-grpc/pb"
	"go-grpc/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer *AuthServer) SignInUser(ctx context.Context, request *pb.SignInUserInput) (*pb.SignInUserResponse, error) {
	user, err := authServer.userService.FindUserByEmail(request.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password")
	}

	if !user.Verified {
		return nil, status.Errorf(codes.PermissionDenied, "user not verified")
	}

	if err := utils.VerifyPassword(user.Password, request.GetPassword()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password")
	}

	access_token, err := utils.CreateToken(authServer.config.AccessTokenExpiresIn, user.ID, authServer.config.AccessTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Could not create access token")
	}

	refresh_token, err := utils.CreateToken(authServer.config.RefreshTokenExpiresIn, user.ID, authServer.config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Could not create refresh token")
	}

	res := &pb.SignInUserResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	return res, nil
}
