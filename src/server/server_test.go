package server

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/random_movies/models"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() *Server {
	return &Server{
		DB:     &MockDB{},
		OMDB:   &MockOMDB{},
		Render: render.New(),
		Log:    logrus.New(),
		Config: &Config{},
	}
}

type MockOMDB struct{}

func (o *MockOMDB) Find(title string) (*models.Movie, error) {
	return testMovie, nil
}

type MockDB struct{}

var testMovie = &models.Movie{}

func (db *MockDB) GetAll() ([]*models.Movie, error) {
	return []*models.Movie{testMovie}, nil
}

func (db *MockDB) GetRandom() (*models.Movie, error) {
	return testMovie, nil
}

func (db *MockDB) Get(imdbID string) (*models.Movie, error) {
	return testMovie, nil
}

func (db *MockDB) Save(movie *models.Movie) error {
	return nil
}

func (db *MockDB) Delete(imdbID string) error {
	return nil
}

func (db *MockDB) MarkSeen(imdbID string) error {
	return nil
}

func TestRandomMovie(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	s := newTestServer()
	s.getRandomMovie(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Random movie did not return %v", http.StatusOK)
	}
}

func TestAddMovie(t *testing.T) {
	jsonStr := []byte(`{"Title":"The Martian"}`)
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	s := newTestServer()
	s.addMovie(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Add movie did not return %v", http.StatusCreated)
	}
}
