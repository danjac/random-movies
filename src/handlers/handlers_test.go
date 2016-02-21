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

type fakeRepo struct {
	movies []*models.Movie
	movie  *models.Movie
	err    error
}

func (r *fakeRepo) GetAll() ([]*models.Movie, error) {
	return r.movies, r.err
}

func (r *fakeRepo) GetRandom() (*models.Movie, error) {
	return r.movie, r.err
}

func (r *fakeRepo) Get(_ string) (*models.Movie, error) {
	return r.movie, r.err
}

func (r *fakeRepo) Delete(_ string) error      { return r.err }
func (r *fakeRepo) MarkSeen(_ string) error    { return r.err }
func (r *fakeRepo) Save(_ *models.Movie) error { return r.err }

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
		Repo:    &fakeRepo{movie: movie},
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
		Repo:    &fakeRepo{err: store.ErrMovieNotFound},
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
		Repo:    &fakeRepo{movie: movie},
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
