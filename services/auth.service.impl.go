package services

import (
	"context"
	"errors"
	"go-grpc/config"
	"go-grpc/models"
	"go-grpc/utils"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(collection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImpl{
		collection, ctx,
	}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := uc.collection.InsertOne(uc.ctx, &user)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exists")
		}
		return nil, err
	}

	// create unique index for email
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("could not create email index")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = uc.collection.FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(user *models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}

func (ac *AuthServiceImpl) ForgotPassword(email string, template *template.Template) error {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}

	err := ac.collection.FindOne(ac.ctx, query).Decode(&user)
	if err != nil {
		log.Printf("[ForgotPassword] error : %v", err)
		if err == mongo.ErrNoDocuments {
			return err
		}
		return err
	}

	if !user.Verified {
		log.Printf("[ForgotPassword] user %s is unverified and send request forgot password", user.Email)
		return errors.New("user is not verified")
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("[ForgotPassword] could not load config: %v", err)
		return err
	}

	resetToken := randstr.String(20)
	filter := bson.D{{Key: "email", Value: user.Email}}
	data := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "password_reset_token", Value: utils.Encode(resetToken)},
			{Key: "password_reset_at", Value: time.Now().Add(time.Minute * 15)},
		},
	}}

	_, err = ac.collection.UpdateOne(ac.ctx, filter, data)
	if err != nil {
		log.Printf("[ForgotPassword] error when updaate collection: %v", err)
		return err
	}

	var firstName = user.Name
	emailData := utils.EmailData{
		URL:       config.Origin + "/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	_ = utils.SendEmail(user, &emailData, template, "resetPassword.html")

	return nil
}

func (ac *AuthServiceImpl) ResetPassword(token string, input *models.ResetPasswordInput) error {
	var user *models.DBResponse

	filter := bson.D{{Key: "password_reset_token", Value: token}}
	err := ac.collection.FindOne(ac.ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[ResetPassword] error no user found with token : %s", token)
		} else {
			log.Printf("[ResetPassword] operation error when get document with token : %s", token)
		}
		return err
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("[ResetPassword] error when load config : %v", err)
		return err
	}

	// in minute
	resetPasswordExpiredTime := time.Duration(config.ResetPasswordExpireTime)

	if time.Now().After(user.PasswordResetAt.Add(time.Minute * resetPasswordExpiredTime)) {
		log.Println("[ResetPassword] reset password expired")
		return errors.New("token expired")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Printf("[ResetPassword] error when hash password : %v", err)
		return err
	}

	filter = bson.D{{Key: "email", Value: user.Email}}
	updateQuery := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "password",
					Value: hashedPassword,
				},
				{
					Key:   "password_reset_token",
					Value: nil,
				},
				{
					Key:   "password_reset_at",
					Value: nil,
				},
			},
		},
	}
	_, err = ac.collection.UpdateOne(ac.ctx, filter, updateQuery)
	if err != nil {
		log.Printf("[ResetPassword] error: user %s when updating new password", user.Email)
		return err
	}

	return nil
}
