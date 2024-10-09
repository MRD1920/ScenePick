package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	elasticSearch "github.com/mrd1920/ScenePick/src/services/elastic_search"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Server) transferMoviesToElasticSearch(ctx *gin.Context) {
	// Define a slice to hold the results
	var movies []bson.M

	// Define the filter (empty filter to fetch all records)
	filter := bson.M{}

	// Define the options (optional)
	findOptions := options.Find().SetBatchSize(5000)

	// Execute the query
	// cursor, err := s.DBMrg.MongoClient.Database("ScenePick").Collection("Movies").Find(context.TODO(), filter, findOptions)
	cursor, err := s.DBMrg.MongoClient.Database("ScenePick").Collection("Combineddatas").Find(context.TODO(), filter, findOptions)

	if err != nil {
		log.Fatal("error fetching all records from MongoDB ", err)
	}
	defer cursor.Close(s.DBMrg.Ctx)
	log.Println("Fetched all records from MongoDB")

	// Iterate through the cursor and decode each document
	for cursor.Next(s.DBMrg.Ctx) {
		var movie bson.M
		if err := cursor.Decode(&movie); err != nil {
			log.Fatal("Error decoding document ", err)
		}
		// log.Println("Movie: ", movie)
		movies = append(movies, movie)
	}

	log.Println("Total Movies fetched from Database: ", len(movies))

	// Check for errors during iteration
	// if err := cursor.Err(); err != nil {
	// 	log.Fatal("Error while iterating ", err)
	// }

	// Now you have all the movies in the `movies` slice
	// You can proceed to transfer them to Elasticsearch or perform other operations
	// Initialize a counter for the number of movies transferred
	moviesTransferred := 0
	for _, movie := range movies {
		// convertedMovie := elasticSearch.BsonToMovie(movie)
		convertedMovie := elasticSearch.BsonToCombinedMovieCast(movie)

		// elasticSearch.InsertMovies(s.EsClient.Client, []elasticSearch.Movie{convertedMovie})
		elasticSearch.InsertCombinedMovieCast(s.EsClient.Client, []elasticSearch.CombinedMovieCast{convertedMovie})
		moviesTransferred++

	}
	log.Println("Movies transferred to ElasticSearch: ", moviesTransferred)

	ctx.JSON(http.StatusOK, gin.H{"message": "Movies transferred to ElasticSearch!", "Total Movies Transferred": moviesTransferred, "Total movies fetched from MongoDB": len(movies)})

}

func (s *Server) searchMoviesInElasticSearch(ctx *gin.Context) {
	// searchField := []string{"title_x", "overview", "genres.name", "cast.name", "keywords.name", "production_companies.name"}
	searchValue := ctx.Query("value")
	log.Println("Searching for: ", searchValue)

	// movies, err := elasticSearch.QueryMovies(s.EsClient.Client, searchField, searchValue)
	// if err != nil {
	// 	ctx.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	response, err := elasticSearch.QueryMoviesHTTP(searchValue)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, response)
}
