#!/bin/bash

# Step 1: Start the docker-compose services (Elasticsearch and API)
echo "Starting docker-compose services..."
docker-compose up -d

# # Step 2: Wait for Elasticsearch to be healthy
# echo "Waiting for Elasticsearch to be healthy..."
# until [ "$(docker inspect -f {{.State.Health.Status}} elasticsearch)" == "healthy" ]; do
#     sleep 5
#     echo "Waiting for Elasticsearch to be ready..."
# done
# echo "Elasticsearch is healthy!"

# Step 2: Wait for Elasticsearch to be ready by querying the health endpoint directly
echo "Waiting for Elasticsearch to be healthy..."
for i in {1..30}; do  # Try up to 30 times
    if curl --silent --fail http://localhost:9200/_cluster/health; then
        echo "Elasticsearch is healthy!"
        break
    else
        echo "Waiting for Elasticsearch to be ready..."
        sleep 5  # Wait for 5 seconds before retrying
    fi
done

# Check if we exited the loop because we succeeded
if [ $? -ne 0 ]; then
    echo "Elasticsearch did not become healthy in time."
    exit 1  # Exit with an error code
fi

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
