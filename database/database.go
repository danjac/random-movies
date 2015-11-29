package database

import (
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/models"
	"gopkg.in/redis.v3"
)

type MovieGetter interface {
	GetAll() ([]*models.Movie, error)
	GetRandom() (*models.Movie, error)
	Get(string) (*models.Movie, error)
}

type MovieStore interface {
	Delete(string) error
	Save(*models.Movie) error
}

type DB interface {
	MovieGetter
	MovieStore
}

type Config struct {
	URL      string
	Password string
	DB       int64
}

func DefaultConfig() *Config {
	return &Config{
		URL:      "redis:6379",
		Password: "",
		DB:       0,
	}
}

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
		return nil, errors.ErrMovieNotFound
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
			return nil, errors.ErrMovieNotFound
		} else {
			return nil, err
		}
	}
	return movie, nil
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
