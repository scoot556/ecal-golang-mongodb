package handlers

import (
	"ecal-mongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMoviesHandler(db *mongo.Client, c *gin.Context) {
	movieCollection, err := models.GetMovies(db)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(200, gin.H{
		"data": movieCollection,
	})
}

func GetMovieHandler(db *mongo.Client, c *gin.Context) {
	id := c.Param("id")

	c.Header("Content-Type", "application/json")

	result, err := models.GetMovie(db, id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": result,
	})
}
