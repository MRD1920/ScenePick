package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`                               // MongoDB ObjectID
	Name            string             `bson:"name" json:"name" binding:"required:"`        // Name of the user
	Email           string             `bson:"email" json:"email" binding:"required:"`      // Email of the user
	Password        string             `bson:"password" json:"password" binding:"required"` // Password of the user
	FavouriteMovies []int              `bson:"favourite_movies" json:"favourite_movies"`    // List of favourite movie IDs
	WatchLater      []int              `bson:"watch_later" json:"watch_later"`              // List of movie IDs to watch later
	AccessToken     string             `bson:"access_token" json:"access_token"`            // Access token for the user
	RefreshToken    string             `bson:"refresh_token" json:"refresh_token"`          // Refresh token for the user
}
