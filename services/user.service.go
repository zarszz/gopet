package services

import (
	"go-grpc/models"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	SetUserVerificationCode(string, string, string) (*models.DBResponse, error)
}
