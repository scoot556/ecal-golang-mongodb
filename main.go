package main

import (
	"context"
	"ecal-mongo/api/comments"
	"ecal-mongo/api/movies"
	"ecal-mongo/configs"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	router := gin.Default()
	db := configs.ConnectDB()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello from Gin-gonic & mongoDB",
		})
	})

	router.GET("/collections", func(c *gin.Context) {
		collections := GetCollections(db)
		c.JSON(200, gin.H{
			"data": collections,
		})
	})

	router.GET("/movies", func(c *gin.Context) {
		movieCollection, err := movies.GetMovies(db)
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
	})

	router.GET("/movies/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.Header("Content-Type", "application/json")

		result, err := movies.GetMovie(db, id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": result,
		})
	})

	router.GET("/comments", func(c *gin.Context) {
		commentCollection, err := comments.GetComments(db)
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
	})

	router.GET("/comments/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.Header("Content-Type", "application/json")

		result, err := comments.GetComment(db, id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": result,
		})
	})

	router.POST("/comments", func(c *gin.Context) {
		id := c.Query("movieId")
		name := c.Query("name")
		email := c.Query("email")
		comment := c.Query("comment")

		c.Header("Content-Type", "application/json")

		err := comments.CreateComment(db, id, name, email, comment)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": "Comment created successfully",
		})
	})

	router.PATCH("/comments", func(c *gin.Context) {
		id := c.Query("id")
		comment := c.Query("comment")

		c.Header("Content-Type", "application/json")

		err := comments.UpdateComment(db, id, comment)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": "Comment updated successfully",
		})

	})

	router.DELETE("/comments", func(c *gin.Context) {
		id := c.Query("id")

		c.Header("Content-Type", "application/json")

		err := comments.DeleteComment(db, id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"data": "Comment deleted successfully",
		})
	})

	router.Run("localhost:6000")
}

func GetCollections(client *mongo.Client) []string {
	db := client.Database("ecal")

	checkCollections, err := db.ListCollectionNames(context.Background(), options.ListCollections())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(checkCollections)
	return checkCollections
}

func GetCollection(client *mongo.Client, collectionName string) ([]bson.M, error) {
	filter := bson.M{}

	db := client.Database("ecal")

	//findOptions := options.Find().SetLimit(100)

	getCollection, err := db.Collection(collectionName).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer getCollection.Close(context.Background())

	var results []bson.M

	for getCollection.Next(context.Background()) {
		var doc bson.M
		if err := getCollection.Decode(&doc); err != nil {
			return nil, err
		}

		results = append(results, doc)
	}

	if err := getCollection.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
