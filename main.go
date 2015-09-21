package main

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
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

var titles = []string{
	"The Man from UNCLE",
	"Tinker, Tailor, Soldier, Spy",
	"Casino Royale",
	"Jack Reacher",
	"Mission Impossible Rogue Nation",
	"Transporter 2",
}

var titlesMap = make(map[string]bool)

var cache = make(map[string]*Movie)

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

func getMovie(title string) (*Movie, error) {
	if title == "" {
		title = titles[rand.Intn(len(titles))]
	}
	cached, ok := cache[title]
	if ok {
		return cached, nil
	}
	movie, err := getMovieFromOMDB(title)
	if err != nil {
		return nil, err
	}
	cache[title] = movie
	return movie, nil
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func main() {

	for _, title := range titles {
		titlesMap[title] = true
	}

	r := gin.Default()

	r.Use(static.Serve("/static", static.LocalFile("static", true)))
	r.LoadHTMLGlob("templates/*")

	r.GET("/", indexPage)

	api := r.Group("/api/")

	api.GET("/", func(c *gin.Context) {
		title := c.Query("title")
		movie, err := getMovie(title)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, movie)
	})

	api.GET("/titles", func(c *gin.Context) {
		c.JSON(http.StatusOK, titles)
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
			ok, _ := titlesMap[movie.Title]
			if !ok {
				titles = append(titles, movie.Title)
				titlesMap[movie.Title] = true
			}
			cache[s.Title] = movie
		}
		c.JSON(http.StatusOK, movie)
	})

	if err := r.Run(":4000"); err != nil {
		panic(err)
	}

}
