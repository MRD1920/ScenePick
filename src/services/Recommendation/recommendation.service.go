package recommendation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Recommendation struct {
}

// Fetches recommendations for a given movie ID with similar movies from TMDB
func GetRecommendataions(apiKey string, movieId string) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/recommendations?api_key=%s", movieId, apiKey)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("Error creating request: ", err)
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("accept", "application/json")

	httpResponse, err := client.Do(req)
	if err != nil {
		log.Println("Error making request: ", err)
		return nil, err
	}

	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil, err
	}
	return jsonResponse, nil

}
