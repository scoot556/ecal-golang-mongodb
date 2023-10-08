package routes

import (
	"ecal-mongo/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupMovieRoutes(db *mongo.Client, router *gin.Engine) {
	movieGroup := router.Group("/movies")
	{
		movieGroup.GET("/", func(c *gin.Context) {
			handlers.GetMoviesHandler(db, c)
		})
		movieGroup.GET("/:id", func(c *gin.Context) {
			handlers.GetMovieHandler(db, c)
		})
	}
}
