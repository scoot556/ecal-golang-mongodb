package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Comment struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	MovieID primitive.ObjectID `bson:"movie_id"`
	Text    string             `bson:"text"`
	Date    time.Time          `bson:"date"`
}

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrEmptyField      = errors.New("field cannot be empty")
)

//var commentCollection = db.Database("ecal").Collection("comments")

func GetComments(client *mongo.Client) ([]bson.M, error) {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")
	filter := bson.M{}

	getComments, err := commentCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer func(getComments *mongo.Cursor, ctx context.Context) {
		err := getComments.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(getComments, context.Background())

	var results []bson.M

	for getComments.Next(context.Background()) {
		var doc bson.M
		if err := getComments.Decode(&doc); err != nil {
			return nil, err
		}

		results = append(results, doc)
	}

	if err := getComments.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetComment(client *mongo.Client, commentID string) (Comment, error) {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")

	if commentID == "" {
		return Comment{}, errors.New("id cannot be empty")
	}

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return Comment{}, err
	}

	var comment Comment
	checkCommentID := commentCollection.FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		if errors.Is(checkCommentID.Err(), mongo.ErrNoDocuments) {
			return Comment{}, ErrCommentNotFound
		}
		return Comment{}, checkCommentID.Err()
	}

	if err := checkCommentID.Decode(&comment); err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func CreateComment(client *mongo.Client, movieID string, name string, email string, comment string) error {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")
	movieCollection := db.Collection("movies")

	objMovie, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}

	if name == "" || email == "" || comment == "" {
		return ErrEmptyField
	}

	fmt.Println(bson.M{"_id": objMovie})
	checkMovieID := movieCollection.FindOne(context.Background(), bson.M{"_id": objMovie})
	if checkMovieID.Err() != nil {
		if errors.Is(checkMovieID.Err(), mongo.ErrNoDocuments) {
			return errors.New("movie not found")
		}
		return checkMovieID.Err()
	}

	doc := Comment{
		Name:    name,
		Email:   email,
		MovieID: objMovie,
		Text:    comment,
		Date:    time.Now(),
	}

	_, err = commentCollection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateComment(client *mongo.Client, commentID string, comment string) error {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	if comment == "" {
		return ErrEmptyField
	}

	checkCommentID := commentCollection.FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		if errors.Is(checkCommentID.Err(), mongo.ErrNoDocuments) {
			return ErrCommentNotFound
		}
		return checkCommentID.Err()
	} else {
		filter := bson.D{{"_id", objComment}}
		update := bson.D{{"$set", bson.D{{"text", comment}}}}

		_, err = commentCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		return nil
	}
}

func DeleteComment(client *mongo.Client, commentID string) error {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	checkCommentID := commentCollection.FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		if errors.Is(checkCommentID.Err(), mongo.ErrNoDocuments) {
			return ErrCommentNotFound
		}
		return checkCommentID.Err()
	}

	filter := bson.D{{"_id", objComment}}

	_, err = commentCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
