package handlers

import (
	"net/http"
	"time"

	"omdb"
	"store"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"github.com/justinas/nosurf"
	"github.com/turtlemonvh/gin-wraphh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const socketWaitFor = 15 * time.Second

func getConfig(c *gin.Context) *Config {
	cfg, _ := c.Get("cfg")
	return cfg.(*Config)
}

// New returns new Application implementation
func New(Store store.Store, options Options) *Application {

	return &Application{
		Config: &Config{
			Store:      Store,
			OMDB:       omdb.New(),
			Options:    options,
			Downloader: &imdbPosterDownloader{},
		},
	}
}

// Options holds settings and env variables
type Options struct {
	Env,
	StaticURL,
	StaticDir,
	DevServerURL string
	Port int
}

// Config holds all configuration objects
type Config struct {
	Options    Options
	OMDB       omdb.Finder
	Store      store.Store
	Downloader downloader
}

// Application runs everything
type Application struct {
	*Config
}

// Run runs the app
func (app *Application) Run() error {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// inject Config into context

	r.Use(func(c *gin.Context) {
		c.Set("cfg", app.Config)
		c.Next()
	})

	// CSRF
	r.Use(wraphh.WrapHH(nosurf.NewPure))

	r.Static(app.Options.StaticURL, app.Options.StaticDir)

	r.GET("/", indexPage)

	api := r.Group("/api")

	api.GET("/all/", getMovies)
	api.GET("/", getRandomMovie)
	api.POST("/", addMovie)
	api.GET("/suggest", suggest)
	api.GET("/movie/:id", getMovie)
	api.DELETE("/movie/:id", deleteMovie)
	api.PATCH("/seen/:id", markSeen)

	r.Run()
	return nil
}

func indexPage(c *gin.Context) {

	cfg := getConfig(c)

	var staticHost string

	if cfg.Options.Env == "dev" {
		staticHost = cfg.Options.DevServerURL
	}

	csrfToken := nosurf.Token(c.Request)

	data := gin.H{
		"staticHost": staticHost,
		"env":        cfg.Options.Env,
		"csrfToken":  csrfToken,
	}
	c.HTML(http.StatusOK, "index.tmpl", data)
}

func getMovies(c *gin.Context) {
	cfg := getConfig(c)

	movies, err := cfg.Store.GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func markSeen(c *gin.Context) {

	cfg := getConfig(c)

	if err := cfg.Store.MarkSeen(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "Movie seen")
}

func getRandomMovie(c *gin.Context) {

	cfg := getConfig(c)

	movie, err := cfg.Store.GetRandom()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

func suggest(c *gin.Context) {

	cfg := getConfig(c)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for {

		movie, err := cfg.Store.GetRandom()

		if err != nil {
			continue
		}

		if err := conn.WriteJSON(movie); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		time.Sleep(socketWaitFor)
	}

}

func getMovie(c *gin.Context) {
	cfg := getConfig(c)
	movie, err := cfg.Store.Get(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movie)
}

func deleteMovie(c *gin.Context) {
	cfg := getConfig(c)
	if err := cfg.Store.Delete(c.Param("id")); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, "Movie deleted")
}

func addMovie(c *gin.Context) {

	cfg := getConfig(c)

	d := &struct {
		Title string `binding:"required"`
	}{}
	if err := c.Bind(d); err != nil {
		return
	}

	movie, err := cfg.OMDB.Find(d.Title)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	oldMovie, err := cfg.Store.Get(movie.ImdbID)

	if err == store.ErrMovieNotFound {

		if movie.Poster != "" {
			go func(url, imdbID string) {
				if filename, err := cfg.Downloader.download(
					cfg.Options.StaticDir, url, imdbID); err != nil {
					movie.Poster = ""
				} else {
					movie.Poster = filename
				}
				if err := cfg.Store.Save(movie); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
			}(movie.Poster, movie.ImdbID)
		}

		movie.Poster = "" // so we don't get a bad link

		if err := cfg.Store.Save(movie); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, movie)
		return
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, oldMovie)

}
