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
			"data": "Hello from Gin",
		})
	})

	routes.SetupMovieRoutes(db, router)
	routes.SetupCommentRoutes(db, router)

	err := router.Run(":6000")
	if err != nil {
		return
	}
}
