package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	handler "github.com/mrd1920/ScenePick/src/db/Handlers"
	"github.com/mrd1920/ScenePick/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UpdateUserTokens(mongoClient *mongo.Client, email, accessToken, refreshToken string) error {
	//Get the collection
	collection := mongoClient.Database("ScenePick").Collection("users")
	//Update the user with the given email
	filter := bson.M{"email": email}
	//Update the access token and refresh token
	update := bson.M{
		"$set": bson.M{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *Server) login(ctx *gin.Context) {
	var newloginReq loginRequest

	//Bind the incoming request body to the user model
	if err := ctx.ShouldBindBodyWithJSON(&newloginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if the email already exists
	user, err := handler.FindOneUser(s.DBMrg.MongoClient, newloginReq.Email)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			// Send JSON response indicating the user is not found and needs to sign up
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":    "User not found",
				"message":  "Please sign up",
				"redirect": "/register",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Send JSON response indicating the user is found and can login

	// Verify the password
	if !utils.CheckPasswordHash(newloginReq.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.Email, []byte(s.Config.JwtKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Email, []byte(s.Config.JwtKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateUserTokens(s.DBMrg.MongoClient, user.Email, accessToken, refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the access token in a cookie
	ctx.SetCookie("access_token", accessToken, 86400, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 7*86400, "/", "localhost", false, true)

	// Send JSON response indicating successful login and redirect to home page
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"redirect":      "/home",
	})
}
