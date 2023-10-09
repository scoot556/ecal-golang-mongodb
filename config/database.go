package config

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"os"
	"sync"
	"time"
)

var mongoClient *mongo.Client
var mongoInitOnce sync.Once
var mongoInitErr error

func initialiseMongoDB() {
	clientOptions := options.Client().ApplyURI(EnvMongoURI())

	maxRetries := 5
	retryInterval := 10 * time.Second

	var err error

	for retry := 0; retry < maxRetries; retry++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("MongoDB connection attempt %d failed: %v", retry+1, err)
			time.Sleep(retryInterval)
			continue
		}

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Printf("MongoDB ping attempt %d failed: %v", retry+1, err)
			time.Sleep(retryInterval)
			continue
		}

		fmt.Println("Connected to MongoDB")
		return
	}

	log.Fatal("MongoDB initialization failed after maximum retries. Exiting...")
	os.Exit(1)
}
func ConnectDB() *mongo.Client {

	mongoInitOnce.Do(initialiseMongoDB)

	if mongoInitErr != nil {
		log.Fatal("MongoDB initialization failed. Exiting...")
		os.Exit(1)
	}

	return mongoClient
}
