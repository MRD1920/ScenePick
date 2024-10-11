# Movie Recommendation System

## Overview

This is a Movie Recommendation System built using Go (Golang) with the Gin framework. The system allows users to get movie recommendations based on their preferences, manage a watchlist, and utilize Elasticsearch for efficient searching and indexing of movies.

## Features

- **User Authentication**: Secure sign-up and login functionality.
- **Movie Recommendations**: Suggest movies based on user preferences.
- **Watchlist Management**: Add and remove movies from a watchlist.
- **Search Functionality**: Search for movies by name, director, or genre using Elasticsearch.
- **Data Transfer**: Easily transfer movie data to Elasticsearch for better search performance.

## Tech Stack

- **Backend**: Go (Golang) with Gin Framework
- **Database**: MongoDB
- **Search**: Elasticsearch
- **Frontend**: Next.js
- **Containerization**: Docker
- **API Documentation**: Swagger (if applicable)

## Video Explanation

https://www.loom.com/share/7671c03d1b714def9bac55371c3b0db7?sid=cec2cd6f-7b90-43b8-a2cf-5a56073e77d0

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.
- [Docker](https://www.docker.com/products/docker-desktop) and Docker Compose installed.
- [MongoDB](https://www.mongodb.com/try/download/community) installed or access to a MongoDB instance.
- **For Windows Users**: Make sure to have WSL or a Linux environment installed.

### Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/movie-recommendation-system.git
   cd movie-recommendation-system
   ```

2. **Frontend Setup**
   Change to the frontend directory and install the required packages:

   ```bash
   cd frontend
   npm install
   npm run dev
   ```

   This will start the frontend application built with Next.js.

3. **Backend Setup**
   Before running the backend, ensure that your Go modules are tidy:

   ```bash
   go mod tidy
   ```

   change directory to `src`

   ```bash
   cd ./src
   ```

   Once done you can run the `run_app.sh` bash script located inside the `src` directory:

   ```bash
   ./run_app.sh
   ```

   Once set you can now just start the server , you can run the API server directly from the root directory:

   ```bash
   make run
   ```
