#!/bin/bash

# Step 1: Start the docker-compose services (Elasticsearch and API)
echo "Starting docker-compose services..."
docker-compose up -d

# Step 2: Wait for Elasticsearch to be healthy
echo "Waiting for Elasticsearch to be healthy..."
until [ "$(docker inspect -f {{.State.Health.Status}} elasticsearch)" == "healthy" ]; do
    sleep 5
    echo "Waiting for Elasticsearch to be ready..."
done
echo "Elasticsearch is healthy!"

# Step 3: Run the create_es_index.sh script to create the Elasticsearch index
echo "Running the create_es_index.sh script to set up the index..."
./create_es_index.sh

# Step 4: Start the Go API server
echo "Starting the Go API server..."
go run src/main.go &  # Run the API server in the background to allow the next steps

# Step 5: Wait for the API server to be ready to serve requests (basic health check)
echo "Waiting for API server to be ready..."
until $(curl --output /dev/null --silent --head --fail http://localhost:8080); do
    sleep 5
    echo "Waiting for API to be ready..."
done
echo "API server is ready!"

# Step 6: Hit the /api/v1/transfer endpoint to transfer data to Elasticsearch
echo "Hitting the /api/v1/transfer endpoint to transfer data..."
curl -X GET http://localhost:8080/api/v1/transfer

echo "Data transfer to Elasticsearch completed."
