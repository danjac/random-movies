package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Movie struct {
	Title    string
	Actors   string
	Poster   string
	Year     string
	Plot     string
	Director string
	ImdbID   string `json:"imdbID"`
}

func (m *Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Movie) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Movie) Save(db *redis.Client) error {
	return db.Set(m.ImdbID, m, 0).Err()
}

func getMovieFromOMDB(title string) (*Movie, error) {

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

	movie := &Movie{}
	if err := json.Unmarshal(body, movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func getRandomMovie(db *redis.Client) (*Movie, error) {
	imdbID, err := db.RandomKey().Result()
	if err != nil {
		return nil, err
	}
	return getMovie(db, imdbID)
}

func getMovie(db *redis.Client, imdbID string) (*Movie, error) {
	movie := &Movie{}
	if err := db.Get(imdbID).Scan(movie); err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return movie, nil
}

func getMovies(db *redis.Client) ([]*Movie, error) {
	result := db.Keys("*")
	if err := result.Err(); err != nil {
		return nil, err
	}
	var movies []*Movie
	for _, imdbID := range result.Val() {
		movie := &Movie{}
		if err := db.Get(imdbID).Scan(movie); err == nil {
			movies = append(movies, movie)
		}
	}
	return movies, nil
}

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "4000", "server port")
)

func main() {

	flag.Parse()

	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(static.Serve("/static", static.LocalFile("static", true)))
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"env": *env})
	})

	api := r.Group("/api/")

	api.GET("/", func(c *gin.Context) {
		movie, err := getRandomMovie(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, movie)
	})

	api.GET("/movie/:id", func(c *gin.Context) {
		movie, err := getMovie(db, c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if movie == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, movie)
	})

	api.DELETE("/movie/:id", func(c *gin.Context) {
		if err := db.Del(c.Param("id")).Err(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "Deleted")
	})

	api.GET("/all", func(c *gin.Context) {
		movies, err := getMovies(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, movies)
	})

	api.POST("/", func(c *gin.Context) {

		s := &struct {
			Title string `json:"title", binding:"required"`
		}{}

		if err := c.Bind(s); err != nil || s.Title == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		movie, err := getMovieFromOMDB(s.Title)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if movie.Title != "" {
			if err := movie.Save(db); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		c.JSON(http.StatusOK, movie)
	})

	if err := r.Run(":" + *port); err != nil {
		panic(err)
	}

}
