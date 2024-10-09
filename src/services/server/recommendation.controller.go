package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	recommendation "github.com/mrd1920/ScenePick/src/services/Recommendation"
)

func (s *Server) getRecommendations(ctx *gin.Context) {
	//Get the movie ID from the request
	movieId := ctx.Query("movieId")

	//Get the recommendations from the recommendation service
	recommendations, err := recommendation.GetRecommendataions(s.Config.TmdbAPIKey, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	//Return the recommendations
	ctx.JSON(http.StatusOK, recommendations)
}
