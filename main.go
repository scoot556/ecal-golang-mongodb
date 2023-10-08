package main

import (
	"ecal-mongo/config"
	"ecal-mongo/models"
	"ecal-mongo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db := config.ConnectDB()

	models.SetDatabase(db)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello from Gin-gonic & mongoDB",
		})
	})

	routes.SetupMovieRoutes(db, router)
	routes.SetupCommentRoutes(db, router)

	router.Run("localhost:6000")
}
