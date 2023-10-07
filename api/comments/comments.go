package comments

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetComments(client *mongo.Client) ([]bson.M, error) {
	db := client.Database("ecal")
	coll := db.Collection("comments")
	filter := bson.M{}

	getComments, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer getComments.Close(context.Background())

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
	coll := db.Collection("comments")

	if commentID == "" {
		return Comment{}, errors.New("id cannot be empty")
	}

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return Comment{}, err
	}

	var comment Comment
	checkCommentID := coll.FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		if checkCommentID.Err() == mongo.ErrNoDocuments {
			return Comment{}, errors.New("comment not found")
		}
		return Comment{}, err
	}

	if err := checkCommentID.Decode(&comment); err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func CreateComment(client *mongo.Client, movieID string, name string, email string, comment string) error {
	db := client.Database("ecal")

	objMovie, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}

	if name == "" || email == "" || comment == "" {
		return errors.New("name, email, and comment cannot be empty")
	}

	fmt.Println(bson.M{"_id": objMovie})
	checkMovieID := db.Collection("movies").FindOne(context.Background(), bson.M{"_id": objMovie})
	if checkMovieID.Err() != nil {
		return err
	}

	doc := Comment{
		Name:    name,
		Email:   email,
		MovieID: objMovie,
		Text:    comment,
		Date:    time.Now(),
	}

	_, err = db.Collection("comments").InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateComment(client *mongo.Client, commentID string, comment string) error {
	db := client.Database("ecal")
	coll := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	if comment == "" {
		return errors.New("comment cannot be empty")
	}

	checkCommentID := db.Collection("comments").FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		return err
	} else if checkCommentID.Err() == mongo.ErrNoDocuments {
		return err
	} else {
		filter := bson.D{{"_id", objComment}}
		update := bson.D{{"$set", bson.D{{"text", comment}}}}

		_, err = coll.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		return nil
	}
}

func DeleteComment(client *mongo.Client, commentID string) error {
	db := client.Database("ecal")

	coll := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	checkCommentID := db.Collection("comments").FindOne(context.Background(), bson.M{"_id": objComment})
	if checkCommentID.Err() != nil {
		if checkCommentID.Err() == mongo.ErrNoDocuments {
			return errors.New("comment not found")
		}
		return checkCommentID.Err()
	}

	filter := bson.D{{"_id", objComment}}

	_, err = coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
