package movieHandler

import (
	"context"
	"log"
	"time"

	models "github.com/mrd1920/ScenePick/src/db/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertMovie(mongoClient *mongo.Client, movie *models.Movie) (*mongo.InsertOneResult, error) {
	collection := mongoClient.Database("ScenePick").Collection("Movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, movie)

	if err != nil {
		log.Println("Error while inserting movie: ", err)
		return nil, err
	}
	return result, nil
}

func GetMovieById(mongoClient *mongo.Client, id primitive.ObjectID) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movie models.Movie
	err := mongoClient.Database("ScenePick").Collection("Movies").FindOne(ctx, models.Movie{ID: id}).Decode(&movie)

	if err != nil {
		log.Println("Movie not found: ", err)
		return nil, err
	}

	return &movie, nil
}
