package moviedetails

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetMovieDetails(movieID, apiKey string) (map[string]interface{}, error) {
	//New Http client
	client := &http.Client{}
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", movieID, apiKey)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	// req.Header.Add("accept", "application/json")

	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil

}
