package main

import (
	"log"

	"github.com/mrd1920/ScenePick/src/services/server"
	"github.com/mrd1920/ScenePick/src/utils"
)

var Server *server.Server

func main() {
	config, err := utils.LoadConfig("./src")
	if err != nil {
		log.Fatal("Failed to load config")
	}
	Server, err = server.NewServer(config)
	if err != nil {
		log.Fatal("Failed to create server")
	}

	Server.Start(":8080")

}
