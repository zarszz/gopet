package services

import (
	"go-grpc/models"
	"html/template"
)

type AuthService interface {
	SignUpUser(models *models.SignUpInput, template *template.Template) (*models.DBResponse, error)
	SignInUser(models *models.SignInInput) (*string, *string, error)
	ForgotPassword(email string, template *template.Template) error
	ResetPassword(token string, models *models.ResetPasswordInput) error
}
