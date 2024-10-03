package main

import (
	"log"

	"github.com/mrd1920/ScenePick/src/services/server"
	"github.com/mrd1920/ScenePick/src/utils"
)

func main() {
	config, err := utils.LoadConfig("./src")
	if err != nil {
		log.Fatal("Failed to load config")
	}
	server, err := server.NewServer(config)
	if err != nil {
		log.Fatal("Failed to create server")
	}

	server.Start(":8080")

}
