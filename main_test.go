package main

import (
	"bytes"
	"context"
	"ecal-mongo/config"
	"ecal-mongo/models"
	"ecal-mongo/routes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db = config.ConnectDB()

var objectIDString = ""

func TestMain(m *testing.M) {
	models.SetDatabase(db)

	code := m.Run()

	db.Disconnect(context.Background())
	os.Exit(code)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetMoviesHandler(t *testing.T) {
	r := SetUpRouter()
	routes.SetupMovieRoutes(db, r)

	req, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	req, err = http.NewRequest("GET", "/movies/573a1390f29313caabcd4135", nil)
	if err != nil {
		t.Fatal(err)
	}

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func TestGetCommentsHandler(t *testing.T) {
	r := SetUpRouter()
	routes.SetupCommentRoutes(db, r)

	req, err := http.NewRequest("GET", "/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	req, err = http.NewRequest("GET", "/comments/5a9427648b0beebeb69579cf", nil)
	if err != nil {
		t.Fatal(err)
	}

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func TestCreateCommentsHandler(t *testing.T) {
	r := SetUpRouter()
	routes.SetupCommentRoutes(db, r)

	requestBody := models.CreateCommentRequest{
		MovieID: "573a1399f29313caabcec3ce",
		Name:    "John Doe",
		Email:   "johndoe@fakegmail.com",
		Comment: "I really like this movie testing",
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		"POST",
		"/comments",
		bytes.NewBuffer(requestBodyJSON),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusTemporaryRedirect || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("POST", location, bytes.NewBuffer(requestBodyJSON))
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response struct {
		CommentID string `json:"commentID"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	fmt.Println("Object ID", response)

	objectIDString = response.CommentID

	responseBody := w.Body.String()
	fmt.Println(responseBody)
}

func TestUpdateCommentsHandler(t *testing.T) {
	r := SetUpRouter()
	routes.SetupCommentRoutes(db, r)

	fmt.Println("Checking update ID", objectIDString)

	requestBody := models.UpdateCommentRequest{
		CommentID: objectIDString,
		Comment:   "I really hate this movie actually I lied!!! testing",
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		"PATCH",
		"/comments",
		bytes.NewBuffer(requestBodyJSON),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusTemporaryRedirect || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("PATCH", location, bytes.NewBuffer(requestBodyJSON))
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	responseBody := w.Body.String()
	fmt.Println(responseBody)
}

func TestDeleteCommentsHandler(t *testing.T) {
	r := SetUpRouter()
	routes.SetupCommentRoutes(db, r)

	requestBody := models.DeleteCommentRequest{
		CommentID: objectIDString,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		"DELETE",
		"/comments",
		bytes.NewBuffer(requestBodyJSON),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code == http.StatusTemporaryRedirect || w.Code == http.StatusFound {
		location := w.Header().Get("Location")
		req, err := http.NewRequest("DELETE", location, bytes.NewBuffer(requestBodyJSON))
		if err != nil {
			t.Fatal(err)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	responseBody := w.Body.String()
	fmt.Println(responseBody)
}
