package server

import (
	"fmt"
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/decoders"
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/omdb"
	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

type renderer struct {
	templates *template.Template
}

// Render HTML
func (r *renderer) Render(w io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

const SOCKET_WAIT_FOR = 15 * time.Second

func New(db database.DB, config *Config) *Server {
	return &Server{
		DB:     db,
		OMDB:   omdb.New(),
		Config: config,
	}
}

// context globals (not threadsafe, so only store thread-safe objects here)
type Server struct {
	Config *Config
	OMDB   omdb.Finder
	DB     database.DB
}

type Config struct {
	Env,
	StaticURL,
	StaticDir,
	DevServerURL string
	Port int
}

func (s *Server) Router() *http.Server {

	e := echo.New()
	e.SetDebug(true)
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Render HTML

	templates, err := template.ParseGlob(filepath.Join("./templates", "*.tmpl"))
	if err != nil {
		panic(err) // shouldn't do this
	}
	e.SetRenderer(&renderer{templates})

	// static configuration
	e.Static(s.Config.StaticURL, s.Config.StaticDir)

	// index page
	e.Get("/", s.indexPage)

	// API calls
	api := e.Group("/api/")

	api.Get("", s.getRandomMovie)
	api.Post("", s.addMovie)
	api.WebSocket("suggest", s.suggest)
	api.Get("movie/:id", s.getMovie)
	api.Delete("movie/:id", s.deleteMovie)
	api.Patch("seen/:id", s.markSeen)
	api.Get("all/", s.getMovies)

	return e.Server(fmt.Sprintf(":%v", s.Config.Port))
}

func (s *Server) indexPage(c *echo.Context) error {

	var staticHost string

	if s.Config.Env == "dev" {
		staticHost = s.Config.DevServerURL
	}

	csrfToken := nosurf.Token(c.Request())

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        s.Config.Env,
		"csrfToken":  csrfToken,
	}
	return c.Render(http.StatusOK, "index.tmpl", ctx)
}

func (s *Server) markSeen(c *echo.Context) error {
	if err := s.DB.MarkSeen(c.Param("id")); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Movie seen")
}

func (s *Server) getRandomMovie(c *echo.Context) error {

	movie, err := s.DB.GetRandom()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func (s *Server) suggest(c *echo.Context) error {

	ws := c.Socket()
	logger := c.Echo().Logger()

	for {

		for {

			movie, err := s.DB.GetRandom()

			if err != nil {
				logger.Error(err)
				continue
			}

			if err := websocket.JSON.Send(ws, movie); err != nil {
				return err
			}

			time.Sleep(SOCKET_WAIT_FOR)
		}
	}
	return nil

}

func (s *Server) getMovie(c *echo.Context) error {

	movie, err := s.DB.Get(c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movie)
}

func (s *Server) deleteMovie(c *echo.Context) error {
	imdbID := c.Param("id")
	if err := s.DB.Delete(imdbID); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Movie deleted")
}

func (s *Server) getMovies(c *echo.Context) error {

	movies, err := s.DB.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movies)
}

func (s *Server) addMovie(c *echo.Context) error {
	f := &decoders.MovieDecoder{}
	if err := c.Bind(f); err != nil {
		return err
	}

	movie, err := s.OMDB.Find(f.Title)

	if err != nil {
		return err
	}

	oldMovie, err := s.DB.Get(movie.ImdbID)

	if err == errors.ErrMovieNotFound {

		if err := s.DB.Save(movie); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, movie)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, oldMovie)

}
