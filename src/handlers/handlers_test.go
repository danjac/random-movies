package handlers

import (
	"fmt"
	"models"
	"net/http"
	"net/http/httptest"
	"store"
	"testing"

	"github.com/unrolled/render"
)

type fakeDB struct {
	movies []*models.Movie
	movie  *models.Movie
	err    error
}

func (db *fakeDB) GetAll() ([]*models.Movie, error) {
	return db.movies, db.err
}

func (db *fakeDB) GetRandom() (*models.Movie, error) {
	return db.movie, db.err
}

func (db *fakeDB) Get(_ string) (*models.Movie, error) {
	return db.movie, db.err
}

func (db *fakeDB) Delete(_ string) error      { return db.err }
func (db *fakeDB) MarkSeen(_ string) error    { return db.err }
func (db *fakeDB) Save(_ *models.Movie) error { return db.err }

type fakeOMDB struct {
	movie *models.Movie
	err   error
}

func (o *fakeOMDB) Find(_ string) (*models.Movie, error) { return o.movie, o.err }

func TestGetMovie(t *testing.T) {
	movie := &models.Movie{
		ImdbID: "tt090909090",
		Title:  "The Martian",
		Actors: "Matt Damon",
	}
	c := AppConfig{
		DB:      &fakeDB{movie: movie},
		Render:  render.New(),
		Options: Options{},
	}
	router := c.Router()
	r, err := http.NewRequest("GET", fmt.Sprintf("/api/movie/%s", movie.ImdbID), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("Should be a 200 OK")
	}

}

func TestGetMovieNotFound(t *testing.T) {
	c := AppConfig{
		DB:      &fakeDB{err: store.ErrMovieNotFound},
		Render:  render.New(),
		Options: Options{},
	}
	router := c.Router()
	r, err := http.NewRequest("GET", "/api/movie/tt898989", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Error("Should be a 404 Not Found")
	}
}

func TestGetRandomMovie(t *testing.T) {
	movie := &models.Movie{
		ImdbID: "tt090909090",
		Title:  "The Martian",
		Actors: "Matt Damon",
	}
	c := AppConfig{
		DB:      &fakeDB{movie: movie},
		Render:  render.New(),
		Options: Options{},
	}
	router := c.Router()
	r, err := http.NewRequest("GET", "/api/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("Should be a 200 OK")
	}

}
