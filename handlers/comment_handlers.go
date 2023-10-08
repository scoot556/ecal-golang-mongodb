package handlers

import (
	"ecal-mongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCommentsHandler(db *mongo.Client, c *gin.Context) {
	commentCollection, err := models.GetComments(db)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")

	c.JSON(200, gin.H{
		"data": commentCollection,
	})
}

func GetCommentHandler(db *mongo.Client, c *gin.Context) {
	id := c.Param("id")

	c.Header("Content-Type", "application/json")

	result, err := models.GetComment(db, id)
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

func PostCommentHandler(db *mongo.Client, c *gin.Context) {
	id := c.Query("movieId")
	name := c.Query("name")
	email := c.Query("email")
	comment := c.Query("comment")

	c.Header("Content-Type", "application/json")

	err := models.CreateComment(db, id, name, email, comment)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": "Comment created successfully",
	})
}

func UpdateCommentHandler(db *mongo.Client, c *gin.Context) {
	id := c.Query("id")
	comment := c.Query("comment")

	c.Header("Content-Type", "application/json")

	err := models.UpdateComment(db, id, comment)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": "Comment updated successfully",
	})

}

func DeleteCommentHandler(db *mongo.Client, c *gin.Context) {
	id := c.Query("id")

	c.Header("Content-Type", "application/json")

	err := models.DeleteComment(db, id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": "Comment deleted successfully",
	})
}
