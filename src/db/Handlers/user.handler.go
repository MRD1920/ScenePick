package handler

import (
	"context"
	"time"

	models "github.com/mrd1920/ScenePick/src/db/Models"

	"go.mongodb.org/mongo-driver/mongo"
)

func FindOneUser(mongoClient *mongo.Client, email string) (*models.User, error) {
	//Get the collection
	collection := mongoClient.Database("ScenePick").Collection("Users")

	//Find the user with the given email
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := collection.FindOne(ctx, models.User{Email: email})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}

	err := result.Decode((&user))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func InsertOneUser(mongoClient *mongo.Client, newUser models.User) (*mongo.InsertOneResult, error) {
	//Get the collection
	collection := mongoClient.Database("ScenePick").Collection("Users")

	//Insert the new user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return result, nil

}
