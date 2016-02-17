package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"store"
	"time"

	"github.com/asaskevich/govalidator"

	"omdb"

	"github.com/justinas/nosurf"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

func decode(c *echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil {
		return err
	}
	_, err := govalidator.ValidateStruct(data)
	return err
}

type renderer struct {
	templates *template.Template
}

// Render HTML
func (r *renderer) Render(w io.Writer, name string, data interface{}) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

const socketWaitFor = 15 * time.Second

// New returns new server implementation
func New(db store.DB, config *Config) *Server {
	return &Server{
		DB:     db,
		OMDB:   omdb.New(),
		Config: config,
	}
}

// Server is an instance of web app
type Server struct {
	Config *Config
	OMDB   omdb.Finder
	DB     store.DB
}

// Config holds settings and env variables
type Config struct {
	Env,
	StaticURL,
	StaticDir,
	DevServerURL string
	Port int
}

// Run the server instance at given port
func (s *Server) Run() error {

	e := echo.New()
	e.SetDebug(true)
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(nosurf.NewPure)

	// handle not found error
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := h(c)
			if err == store.ErrMovieNotFound {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return err
		}
	})

	// Render HTML

	templates, err := template.ParseGlob(filepath.Join("./templates", "*.tmpl"))
	if err != nil {
		return err
	}
	e.SetRenderer(&renderer{templates})

	// static configuration
	e.Static(s.Config.StaticURL, s.Config.StaticDir)

	//e.Get("/", s.indexPage)

	// API calls
	api := e.Group("/api/")

	api.Get("", s.getRandomMovie)
	api.Post("", s.addMovie)
	api.WebSocket("suggest", s.suggest)
	api.Get("movie/:id", s.getMovie)
	api.Delete("movie/:id", s.deleteMovie)
	api.Patch("seen/:id", s.markSeen)
	api.Get("all/", s.getMovies)

	e.Get("/*", s.indexPage)

	e.Run(fmt.Sprintf(":%v", s.Config.Port))
	return nil

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

			time.Sleep(socketWaitFor)
		}
	}

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
	d := &struct {
		Title string `valid:"required"`
	}{}
	if err := decode(c, d); err != nil {
		return err
	}

	movie, err := s.OMDB.Find(d.Title)

	if err != nil {
		return err
	}

	oldMovie, err := s.DB.Get(movie.ImdbID)

	if err == store.ErrMovieNotFound {

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
