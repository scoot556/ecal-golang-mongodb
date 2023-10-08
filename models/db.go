package models

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client

func SetDatabase(database *mongo.Client) {
	db = database
}
