package gapi

import (
	"context"
	"go-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.GenericResponse, error) {
	verificationCode := req.GetVerificationCode()

	result, err := server.userService.UnsetVerificationCode(verificationCode)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if result == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "could not verify email address")
	}

	res := &pb.GenericResponse{
		Status:  "success",
		Message: "Email verified successfully",
	}

	return res, nil
}
