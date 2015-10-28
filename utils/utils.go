package utils

import (
	"encoding/json"
	"errors"
	"github.com/danjac/random_movies/models"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetMovieFromOMDB(title string) (*models.Movie, error) {

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
		return nil, errors.New("Movie not found")
	}

	return movie, nil
}
