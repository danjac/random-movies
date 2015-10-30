package database

import (
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/models"
	"gopkg.in/redis.v3"
)

type DB struct {
	*redis.Client
}

func (db *DB) GetRandomMovie() (*models.Movie, error) {
	imdbID, err := db.RandomKey().Result()
	if err == redis.Nil {
		return nil, errors.ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return db.GetMovie(imdbID)
}

func (db *DB) GetMovie(imdbID string) (*models.Movie, error) {
	movie := &models.Movie{}
	if err := db.Get(imdbID).Scan(movie); err != nil {
		if err == redis.Nil {
			return nil, errors.ErrMovieNotFound
		} else {
			return nil, err
		}
	}
	return movie, nil
}

func (db *DB) SaveMovie(movie *models.Movie) error {
	return db.Set(movie.ImdbID, movie, 0).Err()
}

func (db *DB) GetMovies() ([]*models.Movie, error) {
	result := db.Keys("*")
	if err := result.Err(); err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for _, imdbID := range result.Val() {
		movie := &models.Movie{}
		if err := db.Get(imdbID).Scan(movie); err == nil {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}
