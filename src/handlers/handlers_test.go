package handlers

import (
	"models"
	"net/http"
	"net/http/httptest"
	"store"
	"testing"

	"github.com/unrolled/render"
)

type mockStore struct {
	mockGetAll    func() ([]*models.Movie, error)
	mockGet       func() (*models.Movie, error)
	mockGetRandom func() (*models.Movie, error)
	mockDelete    func() error
	mockMarkSeen  func() error
	mockSave      func() error
}

func (s *mockStore) GetAll() ([]*models.Movie, error) {
	return s.mockGetAll()
}

func (s *mockStore) GetRandom() (*models.Movie, error) {
	return s.mockGetRandom()
}

func (s *mockStore) Get(_ string) (*models.Movie, error) {
	return s.mockGet()
}

func (s *mockStore) Delete(_ string) error      { return s.mockDelete() }
func (s *mockStore) MarkSeen(_ string) error    { return s.mockMarkSeen() }
func (s *mockStore) Save(_ *models.Movie) error { return s.mockSave() }

type mockOMDB struct {
	mockFind func() (*models.Movie, error)
}

func (o *mockOMDB) Find(_ string) (*models.Movie, error) { return o.mockFind() }

func TestGetMovie(t *testing.T) {

	s := &mockStore{}
	s.mockGet = func() (*models.Movie, error) {
		return &models.Movie{
			ImdbID: "tt090909090",
			Title:  "The Martian",
			Actors: "Matt Damon",
		}, nil
	}

	c := AppConfig{
		Store:   s,
		Render:  render.New(),
		Options: Options{},
	}
	router := c.Router()
	r, err := http.NewRequest("GET", "/api/movie/tt090909090", nil)
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

	s := &mockStore{}
	s.mockGet = func() (*models.Movie, error) {
		return nil, store.ErrMovieNotFound
	}

	c := AppConfig{
		Store:   s,
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

	s := &mockStore{}
	s.mockGetRandom = func() (*models.Movie, error) {
		return &models.Movie{
			ImdbID: "tt090909090",
			Title:  "The Martian",
			Actors: "Matt Damon",
		}, nil
	}
	c := AppConfig{
		Store:   s,
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
