package models

import (
	"context"
	"errors"
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

type CreateCommentRequest struct {
	MovieID string `json:"movieId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

type UpdateCommentRequest struct {
	CommentID string `json:"id"`
	Comment   string `json:"comment"`
}

type DeleteCommentRequest struct {
	CommentID string `json:"id"`
}

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

//func CreateComment(client *mongo.Client, request CreateCommentRequest) error {
//	db := client.Database("ecal")
//	commentCollection := db.Collection("comments")
//	movieCollection := db.Collection("movies")
//
//	objMovie, err := primitive.ObjectIDFromHex(request.MovieID)
//	if err != nil {
//		return err
//	}
//
//	if request.Name == "" || request.Email == "" || request.Comment == "" {
//		return ErrEmptyField
//	}
//
//	movieFilter := bson.M{"_id": objMovie}
//	checkMovieID := movieCollection.FindOne(context.Background(), movieFilter)
//	if checkMovieID.Err() != nil {
//		if errors.Is(checkMovieID.Err(), mongo.ErrNoDocuments) {
//			return errors.New("movie not found")
//		}
//		return checkMovieID.Err()
//	}
//
//	doc := Comment{
//		Name:    request.Name,
//		Email:   request.Email,
//		MovieID: objMovie,
//		Text:    request.Comment,
//		Date:    time.Now(),
//	}
//
//	_, err = commentCollection.InsertOne(context.Background(), doc)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func CreateComment(client *mongo.Client, request CreateCommentRequest) (primitive.ObjectID, error) {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")
	movieCollection := db.Collection("movies")

	objMovie, err := primitive.ObjectIDFromHex(request.MovieID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if request.Name == "" || request.Email == "" || request.Comment == "" {
		return primitive.NilObjectID, ErrEmptyField
	}

	movieFilter := bson.M{"_id": objMovie}
	checkMovieID := movieCollection.FindOne(context.Background(), movieFilter)
	if checkMovieID.Err() != nil {
		if errors.Is(checkMovieID.Err(), mongo.ErrNoDocuments) {
			return primitive.NilObjectID, errors.New("movie not found")
		}
		return primitive.NilObjectID, checkMovieID.Err()
	}

	doc := Comment{
		Name:    request.Name,
		Email:   request.Email,
		MovieID: objMovie,
		Text:    request.Comment,
		Date:    time.Now(),
	}

	result, err := commentCollection.InsertOne(context.Background(), doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateComment(client *mongo.Client, request UpdateCommentRequest) error {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(request.CommentID)
	if err != nil {
		return err
	}

	if request.Comment == "" {
		return ErrEmptyField
	}

	commentFilter := bson.M{"_id": objComment}
	checkCommentID := commentCollection.FindOne(context.Background(), commentFilter)
	if checkCommentID.Err() != nil {
		if errors.Is(checkCommentID.Err(), mongo.ErrNoDocuments) {
			return ErrCommentNotFound
		}
		return checkCommentID.Err()
	} else {
		filter := bson.D{{"_id", objComment}}
		update := bson.D{{"$set", bson.D{{"text", request.Comment}}}}

		_, err = commentCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}

		return nil
	}
}

func DeleteComment(client *mongo.Client, request DeleteCommentRequest) error {
	db := client.Database("ecal")
	commentCollection := db.Collection("comments")

	objComment, err := primitive.ObjectIDFromHex(request.CommentID)
	if err != nil {
		return err
	}

	commentDeleteFilter := bson.M{"_id": objComment}
	checkCommentID := commentCollection.FindOne(context.Background(), commentDeleteFilter)
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
