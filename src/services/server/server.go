package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mrd1920/ScenePick/src/controllers"
	migrate "github.com/mrd1920/ScenePick/src/controllers/Migrate"
	elasticSearch "github.com/mrd1920/ScenePick/src/services/elastic_search"

	DBConfig "github.com/mrd1920/ScenePick/src/db"
	"github.com/mrd1920/ScenePick/src/utils"
)

type Server struct {
	Router   *gin.Engine
	Config   utils.Config
	DBMrg    *DBConfig.DBConfigMgr
	EsClient *elasticSearch.ElasticSearchMgr
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

	esMgr, err := elasticSearch.NewElasticClient(config.ElasticSearchURL)
	if err != nil {
		fmt.Println(err)
		log.Fatalln("Failed to connect to ElasticSearch")
	}
	server.EsClient = esMgr
	return server, nil
}

func (s *Server) SetupRouter() {
	router := gin.Default()
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Routes

	router.GET("/health", controllers.HealthCheck)
	router.GET("/migrate", func(ctx *gin.Context) {
		migrate.Migrate(ctx, s.Config.TmdbAPIKey, s.DBMrg.MongoClient)
	})

	// router.Group("/api/v1/", middleware.AuthMiddleware(s.Config.JwtKey))
	// {
	// 	router.GET("/transfer", s.transferMoviesToElasticSearch)
	// 	router.GET("/essearch", s.searchMoviesInElasticSearch)
	// 	router.GET("/searchmovie", s.searchMovie)
	// 	router.GET("/recommendations", s.getRecommendations)

	// 	//Login routes
	// 	router.POST("/login", s.login)
	// 	router.POST("/signup", s.signup)
	// 	// router.GET("/logout", s.logout)

	// 	//Watchlist and Favourites routes

	// }

	api := router.Group("/api/v1/")
	{
		api.GET("/transfer", s.transferMoviesToElasticSearch)
		api.GET("/essearch", s.searchMoviesInElasticSearch)
		api.GET("/searchmovie", s.searchMovie)
		api.GET("/recommendations", s.getRecommendations)

		//Login routes
		api.POST("/login", s.login)
		api.POST("/signup", s.signup)
		// router.GET("/logout", s.logout)

		//Watchlist and Favourites routes

	}

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
