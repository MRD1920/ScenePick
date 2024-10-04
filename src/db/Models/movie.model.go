package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`     // MongoDB ObjectID
	Adult            bool               `bson:"adult"`             // Whether the movie is for adults
	BackdropPath     string             `bson:"backdrop_path"`     // Path to the backdrop image
	GenreIDs         []int              `bson:"genre_ids"`         // List of genre IDs
	MovieID          int                `bson:"id"`                // TMDB Movie ID
	OriginalLanguage string             `bson:"original_language"` // Original language of the movie
	OriginalTitle    string             `bson:"original_title"`    // Original title of the movie
	Overview         string             `bson:"overview"`          // Short description of the movie
	Popularity       float64            `bson:"popularity"`        // Popularity score
	PosterPath       string             `bson:"poster_path"`       // Path to the poster image
	ReleaseDate      string             `bson:"release_date"`      // Release date of the movie
	Title            string             `bson:"title"`             // Title of the movie
	Video            bool               `bson:"video"`             // Whether the movie has a video
	VoteAverage      float64            `bson:"vote_average"`      // Average vote rating
	VoteCount        int                `bson:"vote_count"`        // Total vote count
}
