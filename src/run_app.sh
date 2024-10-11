#!/bin/bash

GO_EXECUTABLE_PATH="C:\Program Files\Go\bin\go.exe" # Update this path to your Go executable
GO_APP_PATH="./src/main.go"  # Update this path to your Go application
WSL_IP=$(powershell.exe -Command "Get-NetIPAddress -InterfaceAlias 'vEthernet (WSL (Hyper-V firewall))' -AddressFamily IPv4 | Select-Object -ExpandProperty IPAddress" | tr -d '\r\n')
API_URL="http://$WSL_IP:8080/api/v1/transfer"
# Step 1: Start the docker-compose services (Elasticsearch and API)
echo "Starting docker-compose services..."
docker-compose up -d

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

# # Step 4: Start the Go API server
# echo "Starting the Go API server..."
# go run src/main.go &  # Run the API server in the background to allow the next steps

# Step 4: Start the Go API server using Makefile located one level above
echo "\nStarting the Go application using Makefile..."
powershell.exe -Command "make -C .. run" &  # Run make in the parent directory

sleep 10

# Capture the PID of the process using port 8080 using PowerShell
echo "Finding the process ID using port 8080..."
GO_SERVER_PID=$(powershell.exe -Command "Get-NetTCPConnection -LocalPort 8080 | Select-Object -ExpandProperty OwningProcess")
echo "Go server process found with PID $GO_SERVER_PID"

sleep 10 # Wait for the API server to start

#Step 5: Wait for the API server to be ready to serve requests (basic health check)
# echo "Waiting for API server to be ready..."
# until $(curl --output /dev/null --silent --head --fail http://localhost:8080/health); do
#     sleep 5
#     echo "Waiting for API to be ready..."
# done
# echo "API server is ready!"

# Step 6: Hit the /api/v1/transfer endpoint to transfer data to Elasticsearch
echo "Hitting the /api/v1/transfer endpoint to transfer data..."
echo "WSL IP: $WSL_IP"

echo "API URL: $API_URL"
if curl -X GET "http://$WSL_IP:8080/api/v1/transfer"; then
    echo "Data transfer to Elasticsearch completed successfully."
else
    echo "Data transfer to Elasticsearch failed."
fi

# Step 7: Kill the Go server process
if [ -n "$GO_SERVER_PID" ]; then
    echo "Killing the Go server process with PID $GO_SERVER_PID"
    powershell.exe -Command "Stop-Process -Id $GO_SERVER_PID "
else
    echo "No process found using port 8080"
fi