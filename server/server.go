package server

import (
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/logger"
	"github.com/danjac/random_movies/models"
	"github.com/danjac/random_movies/utils"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
)

func NewAppConfig(env string, db *database.DB, staticURL, staticDir, devServerURL string) *AppConfig {
	return &AppConfig{
		DB:           db,
		Env:          env,
		DevServerURL: devServerURL,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		Render:       render.New(),
		Log:          logger.New(),
	}
}

// context globals (not threadsafe, so only store thread-safe objects here)
type AppConfig struct {
	Env          string
	StaticURL    string
	StaticDir    string
	DevServerURL string
	Render       *render.Render
	DB           *database.DB
	Log          *logger.Logger
}

func (c *AppConfig) Abort(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(error).(type) {
	case errors.Error:
		c.Log.Error.Printf("HTTP %s %d: %s", e.Status(), e)
		http.Error(w, e.Error(), e.Status())
	default:
		c.Log.Error.Printf("%v: :%v", r, err)
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
}

func (c *AppConfig) Router() *mux.Router {
	router := mux.NewRouter()

	// static content
	router.PathPrefix(
		c.StaticURL).Handler(http.StripPrefix(
		c.StaticURL, http.FileServer(http.Dir(c.StaticDir))))

	// index page
	router.HandleFunc("/", c.indexPage).Methods("GET")

	// API calls
	api := router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/", c.getRandomMovie).Methods("GET")
	api.HandleFunc("/", c.addMovie).Methods("POST")
	api.HandleFunc("/movie/{id}", c.getMovie).Methods("GET")
	api.HandleFunc("/movie/{id}", c.deleteMovie).Methods("DELETE")
	api.HandleFunc("/all/", c.getMovies).Methods("GET")

	return router
}

func (c *AppConfig) indexPage(w http.ResponseWriter, r *http.Request) {

	var staticHost string

	if c.Env == "dev" {
		staticHost = c.DevServerURL
	}

	csrfToken := nosurf.Token(r)

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        c.Env,
		"csrfToken":  csrfToken,
	}
	c.Render.HTML(w, http.StatusOK, "index", ctx)
}

func (c *AppConfig) getRandomMovie(w http.ResponseWriter, r *http.Request) {

	movie, err := c.DB.GetRandomMovie()
	if err != nil {
		c.Abort(w, r, err)
		return
	}

	if movie == nil {
		c.Abort(w, r, errors.ErrHTTPNotFound)
		return
	}
	c.Render.JSON(w, http.StatusOK, movie)
}

func (c *AppConfig) getMovie(w http.ResponseWriter, r *http.Request) {

	movie, err := c.DB.GetMovie(mux.Vars(r)["id"])
	if err != nil {
		c.Abort(w, r, err)
		return
	}
	if movie == nil {
		c.Abort(w, r, errors.ErrHTTPNotFound)
		return
	}
	c.Render.JSON(w, http.StatusOK, movie)
}

func (c *AppConfig) deleteMovie(w http.ResponseWriter, r *http.Request) {

	if err := c.DB.Del(mux.Vars(r)["id"]).Err(); err != nil {
		c.Abort(w, r, err)
		return
	}
	c.Render.Text(w, http.StatusOK, "Movie deleted")
}

func (c *AppConfig) getMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := c.DB.GetMovies()
	if err != nil {
		c.Abort(w, r, err)
		return
	}
	c.Render.JSON(w, http.StatusOK, movies)
}

func (c *AppConfig) addMovie(w http.ResponseWriter, r *http.Request) {

	f := &models.MovieForm{}
	if err := f.Decode(r); err != nil {
		c.Abort(w, r, errors.HTTPError{http.StatusBadRequest, err})
		return
	}

	movie, err := utils.GetMovieFromOMDB(f.Title)
	if err != nil {
		c.Abort(w, r, err)
		return
	}

	if movie.ImdbID == "" {
		c.Abort(w, r, errors.ErrHTTPNotFound)
		return
	}

	if err := c.DB.SaveMovie(movie); err != nil {
		c.Abort(w, r, err)
		return
	}
	c.Log.Info.Printf("New movie %s added", movie)
	c.Render.JSON(w, http.StatusOK, movie)
}
