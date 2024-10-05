package elasticSearch

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

type Movie struct {
	Adult         bool    `json:"adult"`
	BackdropPath  string  `json:"backdrop_path"`
	GenreIds      []int   `json:"genre_ids"`
	ID            int     `json:"id"`
	OriginalLang  string  `json:"original_language"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	Popularity    float64 `json:"popularity"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	Title         string  `json:"title"`
	Video         bool    `json:"video"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
}

// Initialize a new ElasticSearch client
func NewElasticClient(url string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %v", err)
	}
	return client, nil
}

// InsertMovies inserts a list of movies into the "movies" index
func InsertMovies(client *elastic.Client, movies []Movie) error {
	ctx := context.Background()

	for _, movie := range movies {
		_, err := client.Index().
			Index("movies").
			BodyJson(movie).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("Error indexing movie: %v", err)
		}
	}

	log.Println("Movies added to ElasticSearch!")
	return nil
}

// QueryMovies searches for movies by title, director, or genre
func QueryMovies(client *elastic.Client, searchField string, searchValue string) ([]Movie, error) {
	ctx := context.Background()

	// Perform search query
	query := elastic.NewMatchQuery(searchField, searchValue)
	searchResult, err := client.Search().
		Index("movies").
		Query(query).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error querying ElasticSearch: %v", err)
	}

	// Parse search results
	var movies []Movie
	for _, hit := range searchResult.Hits.Hits {
		var movie Movie
		err := hit.Source.UnmarshalJSON(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling movie: %v", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
