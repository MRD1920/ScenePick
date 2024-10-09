package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) searchMovie(ctx *gin.Context) {
	searchParam := ctx.Query("movieId")
	intSearchParam, err := strconv.Atoi(searchParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	//Get the movie details from the search service
	collection := s.DBMrg.MongoClient.Database("ScenePick").Collection("Combineddatas")
	filter := bson.M{"id": intSearchParam}
	// var result models.CombinedMovieCast
	var result2 map[string]interface{}
	// movieDetail := s.DBMrg.MongoClient.Database("ScenePick").Collection("Combineddatas").FindOne(context.TODO(), models.CombinedMovieCast{ID: intSearchParam})

	// err = collection.FindOne(context.TODO(), filter).Decode(&result)
	err2 := collection.FindOne(context.TODO(), filter).Decode(&result2)
	// var movie models.Movie
	// var movie bson.M

	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Movie not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving movie details"})
		return
	}
	// log.Println("Movie: ", result)
	log.Println("Movie: ", result2)
	ctx.JSON(http.StatusOK, result2)

}
