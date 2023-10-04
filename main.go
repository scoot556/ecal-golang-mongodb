package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

func main() {
	//router := gin.Default()
	//
	//router.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"data": "Hello from Gin-gonic & mongoDB",
	//	})
	//})

	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGOURI")

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err2 := mongo.Connect(context.Background(), clientOptions)
	if err2 != nil {
		log.Fatal(err2)
	}

	err2 = client.Ping(context.Background(), readpref.Primary())
	if err2 != nil {
		log.Fatal(err2)
	}

	db := client.Database("ecal")

	cursor, err2 := db.ListCollections(context.Background(), bson.D{})
	if err2 != nil {
		log.Fatal(err2)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var collInfo bson.M
		if err := cursor.Decode(&collInfo); err != nil {
			log.Fatal(err)
		}

		// Extract the collection name
		collectionName := collInfo["name"].(string)

		fmt.Println("Collection Name:", collectionName)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//err := router.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
