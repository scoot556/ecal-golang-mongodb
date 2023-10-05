package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Hello from Gin-gonic & mongoDB",
		})
	})
	//collections := []string{"movies", "comments"}

	router.GET("/collections", func(c *gin.Context) {
		collections := GetCollections()
		c.JSON(200, gin.H{
			"data": collections,
		})
	})

	router.GET("/collections/movies", func(c *gin.Context) {
		movieCollection, err := GetCollection("movies")
		if err != nil {
			log.Fatal(err)
		}
		//jsonData, err := json.Marshal(movieCollection)
		//if err != nil {
		//	log.Fatal(err)
		//}
		c.Header("Content-Type", "application/json")
		//c.JSON(200, gin.H{
		//	"data": jsonData,
		//})
		c.JSON(200, gin.H{
			"data": movieCollection,
		})
	})

	router.Run("localhost:6000")
}

func GetCollections() []string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//
	mongoURI := os.Getenv("MONGOURI")
	//
	clientOptions := options.Client().ApplyURI(mongoURI)
	//
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("ecal")

	checkCollections, err := db.ListCollectionNames(context.Background(), options.ListCollections())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(checkCollections)
	return checkCollections
}

func GetCollection(collectionName string) ([]bson.M, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGOURI")

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{}

	db := client.Database("ecal")

	findOptions := options.Find().SetLimit(100)

	getCollection, err := db.Collection(collectionName).Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
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
