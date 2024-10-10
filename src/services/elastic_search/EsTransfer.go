package elasticSearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ReleaseDate         time.Time           `json:"release_date"`
	Revenue             interface{}         `json:"revenue"`
	Runtime             float64             `json:"runtime"`
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

type ElasticSearchMgr struct {
	Client *elastic.Client
}

// Initialize a new ElasticSearch client
func NewElasticClient(url string) (*ElasticSearchMgr, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %v", err)
	}
	log.Println("Connected to ElasticSearch!")
	return &ElasticSearchMgr{Client: client}, nil
	// return client, nil
}

// InsertMovies inserts a list of movies into the "movies" index
func InsertCombinedMovieCast(client *elastic.Client, movies []CombinedMovieCast) error {
	ctx := context.Background()

	for _, movie := range movies {
		_, err := client.Index().
			Index("movies").
			BodyJson(movie).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("error indexing movie: %v", err)
		}
	}

	log.Println("Movies added to ElasticSearch!")
	return nil
}

func InsertMovies(client *elastic.Client, movies []Movie) error {
	ctx := context.Background()

	for _, movie := range movies {
		_, err := client.Index().
			Index("movies").
			BodyJson(movie).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("error indexing movie: %v", err)
		}
	}

	log.Println("Movies added to ElasticSearch!")
	return nil
}

// QueryMovies searches for movies by title, director, or genre
func QueryMovies(client *elastic.Client, searchField []string, searchValue string) ([]CombinedMovieCast, error) {
	ctx := context.Background()

	// Perform search query
	// query := elastic.NewMatchQuery(searchField, searchValue)
	query := elastic.NewMultiMatchQuery(searchValue, "title_x", "overview", "genres.name", "cast.name", "keywords.name", "production_companies.name")
	searchResult, err := client.Search().
		Index("movies").
		Query(query).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("title_x")). // Fetch only "title_x" field
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying ElasticSearch: %v", err)
	}

	// Display the number of hits found
	fmt.Printf("Found %d hits\n", searchResult.TotalHits())
	log.Println("Printing search results:- ", searchResult.Hits.Hits)

	// Parse search results
	var movies []CombinedMovieCast
	// var movies []Movie

	for _, hit := range searchResult.Hits.Hits {
		// var movie Movie
		var movie CombinedMovieCast

		err := hit.Source.UnmarshalJSON(hit.Source)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling movie: %v", err)
		}
		movies = append(movies, movie)
	}

	// Print the data after the search query is completed
	for _, movie := range movies {
		fmt.Printf("Movie: %+v\n", movie)
	}

	return movies, nil
}

func BsonToMovie(bsonMovie bson.M) Movie {
	return Movie{
		Adult:        bsonMovie["adult"].(bool),
		BackdropPath: bsonMovie["backdrop_path"].(string),
		GenreIds:     convertToIntSlice(bsonMovie["genre_ids"]),
		// GenreIds:      bsonMovie["genre_ids"].([]int),
		ID:            int(bsonMovie["id"].(int32)),
		OriginalLang:  bsonMovie["original_language"].(string),
		OriginalTitle: bsonMovie["original_title"].(string),
		Overview:      bsonMovie["overview"].(string),
		Popularity:    bsonMovie["popularity"].(float64),
		PosterPath:    bsonMovie["poster_path"].(string),
		ReleaseDate:   bsonMovie["release_date"].(string),
		Title:         bsonMovie["title"].(string),
		Video:         bsonMovie["video"].(bool),
		VoteAverage:   bsonMovie["vote_average"].(float64),
		VoteCount:     int(bsonMovie["vote_count"].(int32)),
	}
}

func convertToIntSlice(input interface{}) []int {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'genre_ids' is not an array of integers")
	}
	intSlice := make([]int, len(interfaceSlice))
	for i, v := range interfaceSlice {
		intSlice[i] = int(v.(int32))
	}
	return intSlice
}

func BsonToCombinedMovieCast(bsonMovie bson.M) CombinedMovieCast {
	return CombinedMovieCast{
		Budget:              int(bsonMovie["budget"].(int32)),
		Genres:              convertToGenreSlice(bsonMovie["genres"]),
		Homepage:            processTagLine(bsonMovie["homepage"]),
		ID:                  int(bsonMovie["id"].(int32)),
		Keywords:            convertToKeywordSlice(bsonMovie["keywords"]),
		OriginalLanguage:    processOriginalTitle(bsonMovie["original_language"]),
		OriginalTitle:       processOriginalTitle(bsonMovie["original_title"]),
		Overview:            bsonMovie["overview"].(string),
		Popularity:          bsonMovie["popularity"].(float64),
		ProductionCompanies: convertToProductionCompanySlice(bsonMovie["production_companies"]),
		ProductionCountries: convertToProductionCountrySlice(bsonMovie["production_countries"]),
		ReleaseDate:         bsonMovie["release_date"].(primitive.DateTime).Time(),
		Revenue:             processRevenue(bsonMovie["revenue"]),
		Runtime:             bsonMovie["runtime"].(float64),
		SpokenLanguages:     convertToLanguageSlice(bsonMovie["spoken_languages"]),
		Status:              bsonMovie["status"].(string),
		Tagline:             processTagLine(bsonMovie["tagline"]),
		TitleX:              processOriginalTitle(bsonMovie["title_x"]),
		VoteAverage:         bsonMovie["vote_average"].(float64),
		VoteCount:           int(bsonMovie["vote_count"].(int32)),
		TitleY:              processOriginalTitle(bsonMovie["title_y"]),
		Cast:                convertToCastSlice(bsonMovie["cast"]),
	}
}

func convertToGenreSlice(input interface{}) []Genre {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'genres' is not an array of Genre")
	}
	genreSlice := make([]Genre, len(interfaceSlice))
	for i, v := range interfaceSlice {
		genreSlice[i] = Genre{
			ID:   int(v.(bson.M)["id"].(int32)),
			Name: v.(bson.M)["name"].(string),
		}
	}
	return genreSlice
}

func convertToKeywordSlice(input interface{}) []Keyword {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'keywords' is not an array of Keyword")
	}
	keywordSlice := make([]Keyword, len(interfaceSlice))
	for i, v := range interfaceSlice {
		keywordSlice[i] = Keyword{
			ID:   int(v.(bson.M)["id"].(int32)),
			Name: v.(bson.M)["name"].(string),
		}
	}
	return keywordSlice
}

func convertToProductionCompanySlice(input interface{}) []ProductionCompany {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'production_companies' is not an array of ProductionCompany")
	}
	productionCompanySlice := make([]ProductionCompany, len(interfaceSlice))
	for i, v := range interfaceSlice {
		productionCompanySlice[i] = ProductionCompany{
			ID:   int(v.(bson.M)["id"].(int32)),
			Name: v.(bson.M)["name"].(string),
		}
	}
	return productionCompanySlice
}

func convertToProductionCountrySlice(input interface{}) []ProductionCountry {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'production_countries' is not an array of ProductionCountry")
	}
	productionCountrySlice := make([]ProductionCountry, len(interfaceSlice))
	for i, v := range interfaceSlice {
		productionCountrySlice[i] = ProductionCountry{
			ISO:  v.(bson.M)["iso_3166_1"].(string),
			Name: v.(bson.M)["name"].(string),
		}
	}
	return productionCountrySlice
}

func convertToLanguageSlice(input interface{}) []Language {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'spoken_languages' is not an array of Language")
	}
	languageSlice := make([]Language, len(interfaceSlice))
	for i, v := range interfaceSlice {
		languageSlice[i] = Language{
			ISO:  v.(bson.M)["iso_639_1"].(string),
			Name: v.(bson.M)["name"].(string),
		}
	}
	return languageSlice
}

func convertToCastSlice(input interface{}) []Cast {
	if input == nil {
		return nil
	}
	interfaceSlice, ok := input.(bson.A)
	if !ok {
		log.Println("Error: 'cast' is not an array of Cast")
	}
	castSlice := make([]Cast, len(interfaceSlice))
	for i, v := range interfaceSlice {
		castSlice[i] = Cast{
			CastID:    int(v.(bson.M)["cast_id"].(int32)),
			Character: v.(bson.M)["character"].(string),
			CreditID:  v.(bson.M)["credit_id"].(string),
			Gender:    int(v.(bson.M)["gender"].(int32)),
			ID:        int(v.(bson.M)["id"].(int32)),
			Name:      v.(bson.M)["name"].(string),
			Order:     int(v.(bson.M)["order"].(int32)),
		}
	}
	return castSlice
}

func processRevenue(revenue interface{}) int64 {
	switch v := revenue.(type) {
	case int32:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v) // If it's a float64, convert to int64 (you may need rounding)
	default:
		return 0
	}
}

func processTagLine(tagline interface{}) string {
	if tagline == nil {
		return ""
	}
	return tagline.(string)
}

func processOriginalTitle(title interface{}) string {
	switch expression := title.(type) {
	case int32:
		return strconv.Itoa(int(expression))
	}
	if title == nil {
		return ""
	}
	return title.(string)
}

//------------------------------

type queryResponseMovieTitle struct {
	Title string `json:"title_x"`
	ID    int    `json:"id"`
	// Add other fields as needed
}

type SearchResult struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source queryResponseMovieTitle `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func QueryMoviesHTTP(queryString string) ([]queryResponseMovieTitle, error) {
	// Define the Elasticsearch URL
	elasticURL := "http://localhost:9200/movies/_search"

	// Construct the query payload
	queryPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  queryString,
				"fields": []string{"title_x", "overview", "genres.name", "cast.name", "keywords.name", "production_companies.name", "crew.name"},
			},
		},
		"_source": []string{"title_x", "id"},
	}

	// Convert the query payload to JSON
	queryJSON, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query payload: %v", err)
	}

	// Make the HTTP request to Elasticsearch
	resp, err := http.Post(elasticURL, "application/json", bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to Elasticsearch: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response body
	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	// Display the number of hits found
	fmt.Printf("Found %d hits\n", searchResult.Hits.Total.Value)

	// Extract the search results
	var movies []queryResponseMovieTitle
	for _, hit := range searchResult.Hits.Hits {
		movies = append(movies, hit.Source)
	}

	// Print the data after the search query is completed
	for _, movie := range movies {
		fmt.Printf("Movie: %+v\n", movie)
	}

	log.Println("Movies: ", movies)

	return movies, nil
}
