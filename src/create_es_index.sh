#!/bin/bash

# Create an index with the required mapping
curl -X PUT "localhost:9200/movies" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "budget": { "type": "integer" },
      "genres": {
        "type": "nested",
        "properties": {
          "id": { "type": "integer" },
          "name": { "type": "text" }
        }
      },
      "homepage": { "type": "text" },
      "id": { "type": "integer" },
      "keywords": {
        "type": "nested",
        "properties": {
          "id": { "type": "integer" },
          "name": { "type": "text" }
        }
      },
      "original_language": { "type": "text" },
      "original_title": { "type": "text" },
      "overview": { "type": "text" },
      "popularity": { "type": "float" },
      "production_companies": {
        "type": "nested",
        "properties": {
          "name": { "type": "text" },
          "id": { "type": "integer" }
        }
      },
      "production_countries": {
        "type": "nested",
        "properties": {
          "iso_3166_1": { "type": "text" },
          "name": { "type": "text" }
        }
      },
      "release_date": { "type": "date" },
      "revenue": { "type": "long" },
      "runtime": { "type": "integer" },
      "spoken_languages": {
        "type": "nested",
        "properties": {
          "iso_639_1": { "type": "text" },
          "name": { "type": "text" }
        }
      },
      "status": { "type": "text" },
      "tagline": { "type": "text" },
      "title_x": { "type": "text" },
      "vote_average": { "type": "float" },
      "vote_count": { "type": "integer" },
      "title_y": { "type": "text" },
      "cast": {
        "type": "nested",
        "properties": {
          "cast_id": { "type": "integer" },
          "character": { "type": "text" },
          "credit_id": { "type": "text" },
          "gender": { "type": "integer" },
          "id": { "type": "integer" },
          "name": { "type": "text" },
          "order": { "type": "integer" }
        }
      }
    }
  }
}'
