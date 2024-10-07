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

type CombinedMovieCast struct {
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	ID                  int                 `json:"id"`
	Keywords            []Keyword           `json:"keywords"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int64               `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []Language          `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	TitleX              string              `json:"title_x"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	TitleY              string              `json:"title_y"`
	Cast                []Cast              `json:"cast"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Keyword struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type ProductionCountry struct {
	ISO  string `json:"iso_3166_1"`
	Name string `json:"name"`
}

type Language struct {
	ISO  string `json:"iso_639_1"`
	Name string `json:"name"`
}

type Cast struct {
	CastID    int    `json:"cast_id"`
	Character string `json:"character"`
	CreditID  string `json:"credit_id"`
	Gender    int    `json:"gender"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
}
