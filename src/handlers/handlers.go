package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"store"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/gommon/log"

	"omdb"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func decode(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	_, err := govalidator.ValidateStruct(data)
	return err
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

// Config holds settings and env variables
type Config struct {
	Env,
	StaticURL,
	StaticDir,
	DevServerURL string
	Port int
}

type context struct {
	*App
	Render *render.Render
}

type handlerFunc func(*context, http.ResponseWriter, *http.Request) error

func (c *context) handler(fn handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(c, w, r); err != nil {
			if err == store.ErrMovieNotFound {
				http.NotFound(w, r)
				return
			}
			log.Error(err)
			http.Error(w, "An error has occurred", http.StatusInternalServerError)
		}
	})
}

// Run the server instance at given port
func (app *App) Run() error {

	c := &context{
		App: app,
		Render: render.New(render.Options{
			DisableHTTPErrorRendering: true,
		}),
	}

	router := mux.NewRouter()
	//router.StrictSlash(true)

	router.PathPrefix(app.Config.StaticURL).Handler(
		http.StripPrefix(app.Config.StaticURL,
			http.FileServer(http.Dir(app.Config.StaticDir))))

	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/", c.handler(getRandomMovie)).Methods("GET")
	api.Handle("/", c.handler(addMovie)).Methods("POST")
	api.Handle("/suggest", c.handler(suggest)).Methods("GET")
	api.Handle("/movie/{id}", c.handler(getMovie)).Methods("GET")
	api.Handle("/movie/{id}", c.handler(deleteMovie)).Methods("DELETE")
	api.Handle("/seen/{id}", c.handler(markSeen)).Methods("PATCH")
	api.Handle("/all/", c.handler(getMovies)).Methods("GET")

	router.Handle("/{path:.*}", c.handler(indexPage)).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(nosurf.New(router))
	n.Run(fmt.Sprintf(":%v", app.Config.Port))
	return nil
}

func indexPage(c *context, w http.ResponseWriter, r *http.Request) error {

	var staticHost string

	if c.Config.Env == "dev" {
		staticHost = c.Config.DevServerURL
	}

	csrfToken := nosurf.Token(r)

	data := map[string]string{
		"staticHost": staticHost,
		"env":        c.Config.Env,
		"csrfToken":  csrfToken,
	}
	return c.Render.HTML(w, http.StatusOK, "index", data)
}

func markSeen(c *context, w http.ResponseWriter, r *http.Request) error {
	if err := c.DB.MarkSeen(mux.Vars(r)["id"]); err != nil {
		return err
	}
	return c.Render.Text(w, http.StatusOK, "Movie seen")
}

func getRandomMovie(c *context, w http.ResponseWriter, r *http.Request) error {

	movie, err := c.DB.GetRandom()

	if err != nil {
		return err
	}

	return c.Render.JSON(w, http.StatusOK, movie)
}

func suggest(c *context, w http.ResponseWriter, r *http.Request) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	for {

		for {

			movie, err := c.DB.GetRandom()

			if err != nil {
				continue
			}

			if err := conn.WriteJSON(movie); err != nil {
				return err
			}

			time.Sleep(socketWaitFor)
		}
	}
	return nil

}

func getMovie(c *context, w http.ResponseWriter, r *http.Request) error {
	movie, err := c.DB.Get(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	return c.Render.JSON(w, http.StatusOK, movie)
}

func deleteMovie(c *context, w http.ResponseWriter, r *http.Request) error {
	imdbID := mux.Vars(r)["id"]
	if err := c.DB.Delete(imdbID); err != nil {
		return err
	}
	return c.Render.Text(w, http.StatusOK, "Movie deleted")
}

func getMovies(c *context, w http.ResponseWriter, r *http.Request) error {
	movies, err := c.DB.GetAll()
	if err != nil {
		return err
	}
	return c.Render.JSON(w, http.StatusOK, movies)
}

func addMovie(c *context, w http.ResponseWriter, r *http.Request) error {
	d := &struct {
		Title string `valid:"required"`
	}{}
	if err := decode(r, d); err != nil {
		return err
	}

	movie, err := c.OMDB.Find(d.Title)

	if err != nil {
		return err
	}

	oldMovie, err := c.DB.Get(movie.ImdbID)

	if err == store.ErrMovieNotFound {

		if err := c.DB.Save(movie); err != nil {
			return err
		}

		return c.Render.JSON(w, http.StatusCreated, movie)
	}

	if err != nil {
		return err
	}

	return c.Render.JSON(w, http.StatusOK, oldMovie)

}
