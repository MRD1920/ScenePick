package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mrd1920/ScenePick/src/controllers"
	migrate "github.com/mrd1920/ScenePick/src/controllers/Migrate"

	DBConfig "github.com/mrd1920/ScenePick/src/db"
	"github.com/mrd1920/ScenePick/src/utils"
)

type Server struct {
	Router *gin.Engine
	Config utils.Config
	DBMrg  *DBConfig.DBConfigMgr
}

func NewServer(config utils.Config) (*Server, error) {
	server := &Server{
		Config: config,
	}

	server.SetupRouter()
	MongoDbMgr, err := DBConfig.ConnectToMongoDB(config.MongoDbConnectionURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB")
	}
	server.DBMrg = MongoDbMgr
	return server, nil
}

func (s *Server) SetupRouter() {
	router := gin.Default()
	router.GET("/health", controllers.HealthCheck)
	router.GET("/migrate", func(ctx *gin.Context) {
		migrate.Migrate(ctx, s.Config.TmdbAPIKey, s.DBMrg.MongoClient)
	})
	s.Router = router
}

func (s *Server) GetRouter() *gin.Engine {
	return s.Router
}

// Runs the server on a specific address and port.
func (s *Server) Start(address string) error {
	return s.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
