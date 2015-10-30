package omdb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/models"
)

// Finds a movie
type Finder interface {
	Find(title string) (*models.Movie, error)
}

// returns default implementation
func New() Finder {
	return &finderImpl{}
}

type finderImpl struct{}

// Search finds a movie from OMDB
func (impl *finderImpl) Find(title string) (*models.Movie, error) {
	u, _ := url.Parse("http://omdbapi.com")
	q := u.Query()
	q.Set("t", title)

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	movie := &models.Movie{}
	if err := json.Unmarshal(body, movie); err != nil {
		return nil, err
	}

	if movie.ImdbID == "" {
		return nil, errors.ErrMovieNotFound
	}

	return movie, nil
}
