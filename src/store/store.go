package store

import (
	"errors"
	"models"

	"gopkg.in/redis.v3"
)

// ErrMovieNotFound is returned if no movie found
var ErrMovieNotFound = errors.New("Movie not found")

// Repo is a data store following the repository pattern
type Repo interface {
	GetAll() ([]*models.Movie, error)
	GetRandom() (*models.Movie, error)
	Get(string) (*models.Movie, error)
	Delete(string) error
	MarkSeen(string) error
	Save(*models.Movie) error
}

// Config holds database configuration info
type Config struct {
	URL      string
	Password string
	DB       int64
}

// DefaultConfig creates config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		URL:      "localhost:6379",
		Password: "",
		DB:       0,
	}
}

type defaultRepo struct {
	*redis.Client
}

// New returns a new database instance
func New(config *Config) (Repo, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       config.DB,
	})
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &defaultRepo{db}, nil
}

func (r *defaultRepo) Delete(imdbID string) error {
	return r.Del(imdbID).Err()
}

func (r *defaultRepo) MarkSeen(imdbID string) error {
	movie, err := r.Get(imdbID)
	if err != nil {
		return err
	}
	movie.Seen = true
	return r.Save(movie)
}

func (r *defaultRepo) Save(movie *models.Movie) error {
	return r.Set(movie.ImdbID, movie, 0).Err()
}

func (r *defaultRepo) GetRandom() (*models.Movie, error) {
	imdbID, err := r.RandomKey().Result()
	if err == redis.Nil {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.Get(imdbID)
}

func (r *defaultRepo) Get(imdbID string) (*models.Movie, error) {
	movie := &models.Movie{}
	if err := r.Client.Get(imdbID).Scan(movie); err != nil {
		if err == redis.Nil {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (r *defaultRepo) GetAll() ([]*models.Movie, error) {
	result := r.Keys("*")
	if err := result.Err(); err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for _, imdbID := range result.Val() {
		movie := &models.Movie{}
		if err := r.Client.Get(imdbID).Scan(movie); err == nil {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}
