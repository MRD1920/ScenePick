package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mrd1920/ScenePick/src/controllers"
	"github.com/mrd1920/ScenePick/src/utils"
)

type Server struct {
	router *gin.Engine
	config utils.Config
}

func NewServer(config utils.Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.GET("/health", controllers.HealthCheck)
	s.router = router
}

// Runs the server on a specific address and port.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
