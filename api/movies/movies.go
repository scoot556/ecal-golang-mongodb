package movies

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Movie struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	Plot             string             `json:"plot" bson:"plot"`
	Genres           []string           `json:"genres" bson:"genres"`
	Runtime          int                `json:"runtime" bson:"runtime"`
	Rated            string             `json:"rated" bson:"rated"`
	Cast             []string           `json:"cast" bson:"cast"`
	NumMflixComments int                `json:"num_mflix_comments" bson:"num_mflix_comments"`
	Title            string             `json:"title" bson:"title"`
	FullPlot         string             `json:"fullplot" bson:"fullplot"`
	Languages        []string           `json:"languages" bson:"languages"`
	Released         time.Time          `json:"released" bson:"released"`
	Directors        []string           `json:"directors" bson:"directors"`
	Writers          []string           `json:"writers" bson:"writers"`
	Awards           Awards             `json:"awards" bson:"awards"`
	LastUpdated      string             `json:"lastupdated" bson:"lastupdated"`
	Year             int                `json:"year" bson:"year"`
	Imdb             Imdb               `json:"imdb" bson:"imdb"`
	Countries        []string           `json:"countries" bson:"countries"`
	Type             string             `json:"type" bson:"type"`
	Tomatoes         Tomatoes           `json:"tomatoes" bson:"tomatoes"`
}

type Awards struct {
	Wins        int    `json:"wins" bson:"wins"`
	Nominations int    `json:"nominations" bson:"nominations"`
	Text        string `json:"text" bson:"text"`
}

type Imdb struct {
	Rating float64 `json:"rating" bson:"rating"`
	Votes  int     `json:"votes" bson:"votes"`
	ID     int     `json:"id" bson:"id"`
}

type Tomatoes struct {
	Viewer      Viewer    `json:"viewer" bson:"viewer"`
	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
}

type Viewer struct {
	Rating     float64 `json:"rating" bson:"rating"`
	NumReviews int     `json:"numReviews" bson:"numReviews"`
	Meter      int     `json:"meter" bson:"meter"`
}

func GetMovies(client *mongo.Client) ([]bson.M, error) {
	db := client.Database("ecal")
	coll := db.Collection("movies")
	filter := bson.M{}

	getMovies, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer getMovies.Close(context.Background())

	var results []bson.M

	for getMovies.Next(context.Background()) {
		var doc bson.M
		if err := getMovies.Decode(&doc); err != nil {
			return nil, err
		}

		results = append(results, doc)
	}

	if err := getMovies.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetMovie(client *mongo.Client, movieID string) (Movie, error) {
	db := client.Database("ecal")
	coll := db.Collection("movies")

	if movieID == "" {
		return Movie{}, errors.New("id cannot be empty")
	}

	objMovie, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return Movie{}, err
	}

	var movie Movie
	checkMovieID := coll.FindOne(context.Background(), bson.M{"_id": objMovie})
	if checkMovieID.Err() != nil {
		if checkMovieID.Err() == mongo.ErrNoDocuments {
			return Movie{}, errors.New("movie not found")
		}
		return Movie{}, err
	}

	if err := checkMovieID.Decode(&movie); err != nil {
		return Movie{}, err
	}

	return movie, nil
}
