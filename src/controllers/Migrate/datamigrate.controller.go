package migrate

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	movieHandler "github.com/mrd1920/ScenePick/src/db/Handlers"
	models "github.com/mrd1920/ScenePick/src/db/Models"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(ctx *gin.Context, apiKey string, mongoClient *mongo.Client) {
	baseURL := "https://api.themoviedb.org/3/discover/movie"
	apiUrl, err := url.Parse(baseURL)

	if err != nil {
		log.Println("Error parsing URL")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	//Create query parameters
	queryParams := url.Values{}
	queryParams.Add("include_adult", "false")
	queryParams.Add("include_video", "false")
	queryParams.Add("language", "en-US")
	queryParams.Add("page", "1")
	queryParams.Add("sort_by", "popularity.desc")
	queryParams.Add("api_key", apiKey)

	// Add query parameters to the URL
	apiUrl.RawQuery = queryParams.Encode()

	//Create a new HTTP client with timeout
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	//Create a new HTTP Get request
	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		log.Fatal("Error creating request", err)
	}

	//Add headers to the request
	// req.Header.Add("Authorization", "Bearer 0d361c9352cfc93982900b4809d66182")
	req.Header.Add("Accept", "application/json")

	//Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making the API request ", err)
	}

	defer resp.Body.Close()

	//Check the response status message code
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Expected status code 200 but got", resp.StatusCode)
	}

	//Read the response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body", err)
	}

	//Print the response
	// fmt.Println(string(body))

	// Unmarshal the response body into a map
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Fatal("Error unmarshaling response body", err)
	}

	// var movie models.Movie
	for _, result := range jsonResponse["results"].([]interface{}) {
		resultMap := result.(map[string]interface{})

		// Convert genre_ids from []interface{} to []int
		var genreIDs []int
		for _, v := range resultMap["genre_ids"].([]interface{}) {
			if id, ok := v.(float64); ok {
				genreIDs = append(genreIDs, int(id))
			} else {
				log.Println("Error converting genre_id to int")
			}
		}

		movie := models.Movie{
			Adult:            resultMap["adult"].(bool),
			BackdropPath:     resultMap["backdrop_path"].(string),
			GenreIDs:         genreIDs,
			MovieID:          int(resultMap["id"].(float64)),
			OriginalLanguage: resultMap["original_language"].(string),
			OriginalTitle:    resultMap["original_title"].(string),
			Overview:         resultMap["overview"].(string),
			Popularity:       resultMap["popularity"].(float64),
			PosterPath:       resultMap["poster_path"].(string),
			ReleaseDate:      resultMap["release_date"].(string),
			Title:            resultMap["title"].(string),
			Video:            resultMap["video"].(bool),
			VoteAverage:      resultMap["vote_average"].(float64),
			VoteCount:        int(resultMap["vote_count"].(float64)),
		}
		movieHandler.InsertMovie(mongoClient, &movie)
	}
	ctx.JSON(http.StatusOK, jsonResponse)

}
