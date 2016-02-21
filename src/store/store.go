package store

import (
	"errors"
	"models"

	"gopkg.in/redis.v3"
)

// ErrMovieNotFound is returned if no movie found
var ErrMovieNotFound = errors.New("Movie not found")

// Store is a data store following the repository pattern
type Store interface {
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

type defaultStore struct {
	*redis.Client
}

// New returns a new database instance
func New(config *Config) (Store, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       config.DB,
	})
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &defaultStore{db}, nil
}

func (s *defaultStore) Delete(imdbID string) error {
	return s.Del(imdbID).Err()
}

func (s *defaultStore) MarkSeen(imdbID string) error {
	movie, err := s.Get(imdbID)
	if err != nil {
		return err
	}
	movie.Seen = true
	return s.Save(movie)
}

func (s *defaultStore) Save(movie *models.Movie) error {
	return s.Set(movie.ImdbID, movie, 0).Err()
}

func (s *defaultStore) GetRandom() (*models.Movie, error) {
	imdbID, err := s.RandomKey().Result()
	if err == redis.Nil {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.Get(imdbID)
}

func (s *defaultStore) Get(imdbID string) (*models.Movie, error) {
	movie := &models.Movie{}
	if err := s.Client.Get(imdbID).Scan(movie); err != nil {
		if err == redis.Nil {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *defaultStore) GetAll() ([]*models.Movie, error) {
	result := s.Keys("*")
	if err := result.Err(); err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for _, imdbID := range result.Val() {
		movie := &models.Movie{}
		if err := s.Client.Get(imdbID).Scan(movie); err == nil {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}
