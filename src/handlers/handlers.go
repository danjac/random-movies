package handlers

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

const appContextKey = "app"

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
func New(db store.DB, config *Config) *App {
	return &App{
		DB:     db,
		OMDB:   omdb.New(),
		Config: config,
	}
}

// App is an instance of web app
type App struct {
	Config *Config
	OMDB   omdb.Finder
	DB     store.DB
}

func getApp(c *echo.Context) *App {
	return c.Get(appContextKey).(*App)
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
func (app *App) Run() error {

	e := echo.New()
	e.SetDebug(true)
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(nosurf.NewPure)

	// add instance to context
	e.Use(func(c *echo.Context) error {
		c.Set(appContextKey, app)
		return nil
	})

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
	e.Static(app.Config.StaticURL, app.Config.StaticDir)

	// API calls
	api := e.Group("/api/")

	api.Get("", getRandomMovie)
	api.Post("", addMovie)
	api.WebSocket("suggest", suggest)
	api.Get("movie/:id", getMovie)
	api.Delete("movie/:id", deleteMovie)
	api.Patch("seen/:id", markSeen)
	api.Get("all/", getMovies)

	e.Get("/*", indexPage)

	e.Run(fmt.Sprintf(":%v", app.Config.Port))
	return nil

}

func indexPage(c *echo.Context) error {
	app := getApp(c)

	var staticHost string

	if app.Config.Env == "dev" {
		staticHost = app.Config.DevServerURL
	}

	csrfToken := nosurf.Token(c.Request())

	data := map[string]string{
		"staticHost": staticHost,
		"env":        app.Config.Env,
		"csrfToken":  csrfToken,
	}
	return c.Render(http.StatusOK, "index.tmpl", data)
}

func markSeen(c *echo.Context) error {
	if err := getApp(c).DB.MarkSeen(c.Param("id")); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Movie seen")
}

func getRandomMovie(c *echo.Context) error {

	movie, err := getApp(c).DB.GetRandom()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

func suggest(c *echo.Context) error {

	ws := c.Socket()
	logger := c.Echo().Logger()

	for {

		for {

			movie, err := getApp(c).DB.GetRandom()

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

func getMovie(c *echo.Context) error {
	movie, err := getApp(c).DB.Get(c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movie)
}

func deleteMovie(c *echo.Context) error {
	imdbID := c.Param("id")
	if err := getApp(c).DB.Delete(imdbID); err != nil {
		return err
	}
	return c.String(http.StatusOK, "Movie deleted")
}

func getMovies(c *echo.Context) error {
	movies, err := getApp(c).DB.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, movies)
}

func addMovie(c *echo.Context) error {
	d := &struct {
		Title string `valid:"required"`
	}{}
	if err := decode(c, d); err != nil {
		return err
	}

	app := getApp(c)

	movie, err := app.OMDB.Find(d.Title)

	if err != nil {
		return err
	}

	oldMovie, err := app.DB.Get(movie.ImdbID)

	if err == store.ErrMovieNotFound {

		if err := app.DB.Save(movie); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, movie)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, oldMovie)

}
