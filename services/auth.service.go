package services

import (
	"go-grpc/models"
	"html/template"
)

type AuthService interface {
	SignUpUser(models *models.SignUpInput) (*models.DBResponse, error)
	SignInUser(models *models.SignInInput) (*models.DBResponse, error)
	ForgotPassword(email string, template *template.Template) error
	ResetPassword(token string, models *models.ResetPasswordInput) error
}
