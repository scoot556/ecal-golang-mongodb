package setup

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

func main() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGOURI")

	//commentsJSONPath := "../comments.json"
	//moviesJSONPath := "../movies.json"

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return "SETUP"

	//commentsJSON, err := ioutil.ReadFile(commentsJSONPath)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//moviesJSON, err := ioutil.ReadFile(moviesJSONPath)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//dbName := "ecal"

	//if client.Database(dbName).Collection("comments")
}
