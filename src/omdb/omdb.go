package omdb

import (
	"encoding/json"
	"net/http"
	"net/url"

	"httperrors"
	"models"
)

// Finder finds a movie based on title
type Finder interface {
	Find(title string) (*models.Movie, error)
}

// New returns default implementation
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

	movie := &models.Movie{}
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	if movie.ImdbID == "" {
		return nil, httperrors.ErrMovieNotFound
	}

	return movie, nil
}
