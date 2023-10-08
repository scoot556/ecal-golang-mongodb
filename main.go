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

	routes.SetupMovieRoutes(db, router)
	routes.SetupCommentRoutes(db, router)

	router.Run("localhost:6000")
}
