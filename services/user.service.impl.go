package services

import (
	"context"
	"fmt"
	"go-grpc/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (uc *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := uc.collection.FindOne(uc.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}
func (uc *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": strings.ToLower(email)}
	err := uc.collection.FindOne(uc.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserServiceImpl) SetUserVerificationCode(id string, field string, value string) (*models.DBResponse, error) {
	userId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	fmt.Print(result.ModifiedCount)
	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}
