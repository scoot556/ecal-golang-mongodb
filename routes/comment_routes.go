package routes

import (
	"ecal-mongo/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupCommentRoutes(db *mongo.Client, router *gin.Engine) {
	commentGroup := router.Group("/comments")
	{
		//commentGroup.GET("/", handlers.GetCommentsHandler(db))
		//commentGroup.GET("/:id", handlers.GetCommentHandler(db))
		//commentGroup.POST("/", handlers.PostCommentHandler(db))
		//commentGroup.PATCH("/", handlers.UpdateCommentHandler(db))
		//commentGroup.DELETE("/", handlers.DeleteCommentHandler(db))
		commentGroup.GET("/", func(c *gin.Context) {
			handlers.GetCommentsHandler(db, c)
		})
		commentGroup.GET("/:id", func(c *gin.Context) {
			handlers.GetCommentHandler(db, c)
		})
		commentGroup.POST("/", func(c *gin.Context) {
			handlers.PostCommentHandler(db, c)
		})
		commentGroup.PATCH("/", func(c *gin.Context) {
			handlers.UpdateCommentHandler(db, c)
		})
		commentGroup.DELETE("/", func(c *gin.Context) {
			handlers.DeleteCommentHandler(db, c)
		})
	}
}
