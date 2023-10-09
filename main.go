package main

import (
	"ecal-mongo/config"
	"ecal-mongo/models"
	"ecal-mongo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
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

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
