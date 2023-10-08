package handlers

import (
	"ecal-mongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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
	var request models.CreateCommentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")

	commentID, err := models.CreateComment(db, request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data":      "Comment created successfully",
		"commentID": commentID.Hex(),
	})
}

func UpdateCommentHandler(db *mongo.Client, c *gin.Context) {
	var request models.UpdateCommentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")

	err := models.UpdateComment(db, request)
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
	var request models.DeleteCommentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")

	err := models.DeleteComment(db, request)
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
