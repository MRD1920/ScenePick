package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	handler "github.com/mrd1920/ScenePick/src/db/Handlers"
	models "github.com/mrd1920/ScenePick/src/db/Models"
	"github.com/mrd1920/ScenePick/src/utils"
)

func (s *Server) signup(ctx *gin.Context) {
	var user models.User

	// Bind the incoming request body to the user model
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email already exists
	_, err := handler.FindOneUser(s.DBMrg.MongoClient, user.Email)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			// Hash the password
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Update the password with the hashed password
			user.Password = hashedPassword

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

			// Store the access token and refresh token in the user model
			user.AccessToken = accessToken
			user.RefreshToken = refreshToken

			// Insert the new user
			_, err = handler.InsertOneUser(s.DBMrg.MongoClient, user)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Set the access token and refresh token in cookies
			ctx.SetCookie("access_token", accessToken, 86400, "/", "localhost", false, true)
			ctx.SetCookie("refresh_token", refreshToken, 7*86400, "/", "localhost", false, true)

			// Redirect the user to the home page
			ctx.Redirect(http.StatusFound, "/home")
			return
		}
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
}
