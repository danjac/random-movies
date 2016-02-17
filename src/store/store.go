package store

import (
	"errors"
	"models"

	"gopkg.in/redis.v3"
)

// ErrMovieNotFound is returned if no movie found
var ErrMovieNotFound = errors.New("Movie not found")

// MovieReader reads data from the store
type MovieReader interface {
	GetAll() ([]*models.Movie, error)
	GetRandom() (*models.Movie, error)
	Get(string) (*models.Movie, error)
}

// MovieWriter writes data to the store
type MovieWriter interface {
	Delete(string) error
	MarkSeen(string) error
	Save(*models.Movie) error
}

// DB handles reads/writes to/from the store
type DB interface {
	MovieReader
	MovieWriter
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

// New returns a new database instance
func New(config *Config) (DB, error) {
	db := &defaultImpl{redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       config.DB,
	})}
	_, err := db.Ping().Result()
	return db, err
}

type defaultImpl struct {
	*redis.Client
}

func (db *defaultImpl) Delete(imdbID string) error {
	return db.Del(imdbID).Err()
}

func (db *defaultImpl) GetRandom() (*models.Movie, error) {
	imdbID, err := db.RandomKey().Result()
	if err == redis.Nil {
		return nil, ErrMovieNotFound
	}
	if err != nil {
		return nil, err
	}
	return db.Get(imdbID)
}

func (db *defaultImpl) Get(imdbID string) (*models.Movie, error) {
	movie := &models.Movie{}
	if err := db.Client.Get(imdbID).Scan(movie); err != nil {
		if err == redis.Nil {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (db *defaultImpl) MarkSeen(imdbID string) error {
	movie, err := db.Get(imdbID)
	if err != nil {
		return err
	}
	movie.Seen = true
	return db.Save(movie)
}

func (db *defaultImpl) Save(movie *models.Movie) error {
	return db.Set(movie.ImdbID, movie, 0).Err()
}

func (db *defaultImpl) GetAll() ([]*models.Movie, error) {
	result := db.Keys("*")
	if err := result.Err(); err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for _, imdbID := range result.Val() {
		movie := &models.Movie{}
		if err := db.Client.Get(imdbID).Scan(movie); err == nil {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}
