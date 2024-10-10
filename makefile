# Makefile for Go Application

# Define the Go binary
GO_BIN=go

# Define the main application source file
SRC=src/main.go



# Default target to run the application
run: 
	@echo "Building and running the application..."
	$(GO_BIN) run $(SRC)

# Target to build the application
build:
	@echo "Building the application..."
	$(GO_BIN) build -o myapp $(SRC)

# Target to clean up build files
clean:
	@echo "Cleaning up..."
	rm -f myapp
# Define the target for running the application
.PHONY: run build clean