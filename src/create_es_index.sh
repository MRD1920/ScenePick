#!/bin/bash

# Create movies index with appropriate mappings for movie name, director, and genre
curl -X PUT "http://localhost:9200/movies" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "adult": { "type": "boolean" },
      "backdrop_path": { "type": "text" },
      "genre_ids": { "type": "integer" },
      "id": { "type": "integer" },
      "original_language": { "type": "text" },
      "original_title": { "type": "text" },
      "overview": { "type": "text" },
      "popularity": { "type": "float" },
      "poster_path": { "type": "text" },
      "release_date": { "type": "date" },
      "title": { "type": "text" },
      "video": { "type": "boolean" },
      "vote_average": { "type": "float" },
      "vote_count": { "type": "integer" }
    }
  }
}'
