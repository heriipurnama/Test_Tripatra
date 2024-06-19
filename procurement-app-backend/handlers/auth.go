package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"procurement-app-backend/models"
	"procurement-app-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	userCollection = client.Database("procurementdb").Collection("users")
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func Register(ctx context.Context, name string, email string, password string) (*models.User, error) {
	hashedPassword := hashPassword(password)

	user := &models.User{
		UserID:   primitive.NewObjectID().Hex(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Login(ctx context.Context, email string, password string) (*models.LoginResponse, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Password != hashPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Email, user.Name, user.UserID)
	if err != nil {
		return nil, err
	}

	response := &models.LoginResponse{
		Token: token,
		User: models.User{
			UserID: user.UserID, // Assuming the user ID is an ObjectID
			Name:   user.Name,
			Email:  user.Email,
			Role:   user.Role,
		},
	}

	return response, nil
}
