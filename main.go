package main

import (
	"context"
	"ecal-mongo/configs"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
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
		movieCollection, err := GetCollection(db, "movies")
		if err != nil {
			log.Fatal(err)
		}
		c.Header("Content-Type", "application/json")

		c.JSON(200, gin.H{
			"data": movieCollection,
		})
	})

	router.GET("/comments", func(c *gin.Context) {
		commentCollection, err := GetCollection(db, "comments")
		if err != nil {
			log.Fatal(err)
		}

		c.Header("Content-Type", "application/json")

		c.JSON(200, gin.H{
			"data": commentCollection,
		})
	})

	router.POST("/comments", func(c *gin.Context) {
		id := c.Query("id")
		name := c.Query("name")
		email := c.Query("email")
		comment := c.Query("comment")

		CreateComment(db, id, name, email, comment)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//c.JSON(200, gin.H{
		//	"data": createComment,
		//})
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

type Comment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	Text    string             `bson:"text"`
	Date    time.Time          `bson:"date"`
}

func CreateComment(client *mongo.Client, movieID string, name string, email string, comment string) {
	fmt.Println(movieID, name, email, comment)

	db := client.Database("ecal")

	objMovie, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bson.M{"_id": objMovie})
	checkMovieID := db.Collection("movies").FindOne(context.Background(), bson.M{"_id": objMovie})
	if checkMovieID.Err() == nil {
		fmt.Println("Movie ID exists: true")
	} else if checkMovieID.Err() == mongo.ErrNoDocuments {
		fmt.Println("Movie ID exists: false")
	} else {
		fmt.Println("Error checking Movie ID:", checkMovieID.Err())
	}

	doc := Comment{
		Name:    name,
		Email:   email,
		MovieID: objMovie,
		Text:    comment,
		Date:    time.Now(),
	}

	insertComment, err := db.Collection("comments").InsertOne(context.Background(), doc)
	if err != nil {
		fmt.Println("Error inserting document:", err)
	}
	//return name, nil
	fmt.Println("Inserted doc", insertComment)
}

func UpdateComment(commentID string, movieID string) {

}

func DeleteComment(commentID string, movieID string) {

}
