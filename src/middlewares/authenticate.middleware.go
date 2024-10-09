package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtkey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Retreive the access_token from the cookie
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, access token not found"})
			ctx.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtkey), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Token is valid, proceed with the request
		ctx.Next()

	}
}
