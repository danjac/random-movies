package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"goji.io/pat"

	"github.com/asaskevich/govalidator"

	"omdb"
	"store"

	"goji.io"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

const requestContextKey = "reqctx"

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

// New returns new Application implementation
func New(Store store.Store, options Options) *Application {

	return &Application{
		Store:      Store,
		OMDB:       omdb.New(),
		Options:    options,
		Render:     render.New(),
		downloader: &imdbPosterDownloader{},
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

// Application is an instance of web app
type Application struct {
	Options    Options
	OMDB       omdb.Finder
	Store      store.Store
	Render     *render.Render
	downloader downloader
}

// RequestContext wraps Application and provides extra per-request info
// In a larger app we'd keep user identity etc here
type RequestContext struct {
	*Application
	Err error
}

func getRequestContext(ctx context.Context) *RequestContext {
	return ctx.Value(requestContextKey).(*RequestContext)
}

// NewRequestContext creates a new request context
func (app *Application) NewRequestContext(ctx context.Context) (context.Context, *RequestContext) {
	reqctx := &RequestContext{Application: app}
	ctx = context.WithValue(ctx, requestContextKey, reqctx)
	return ctx, reqctx
}

// Router creates http handler
func (app *Application) Router() http.Handler {

	router := goji.NewMux()

	router.Handle(
		pat.Get(app.Options.StaticURL+"*"),
		http.StripPrefix(app.Options.StaticURL,
			http.FileServer(http.Dir(app.Options.StaticDir))))

	api := goji.SubMux()

	api.HandleFuncC(pat.Get("/all/"), getMovies)
	api.HandleFuncC(pat.Get("/"), getRandomMovie)
	api.HandleFuncC(pat.Post("/"), addMovie)
	api.HandleFuncC(pat.Get("/suggest"), suggest)
	api.HandleFuncC(pat.Get("/movie/:id"), getMovie)
	api.HandleFuncC(pat.Delete("/movie/:id"), deleteMovie)
	api.HandleFuncC(pat.Patch("/seen/:id"), markSeen)

	router.HandleC(pat.New("/api/*"), api)

	router.HandleFuncC(pat.Get("/*"), indexPage)

	// inject Application
	router.UseC(func(h goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			// append request context
			ctx, reqctx := app.NewRequestContext(ctx)
			h.ServeHTTPC(ctx, w, r)
			if reqctx.Err != nil {
				if reqctx.Err == store.ErrMovieNotFound {
					http.NotFound(w, r)
				} else {
					log.Println(reqctx.Err)
					http.Error(w, "An error has occurred", http.StatusInternalServerError)
				}
			}

			// check for errors here
		})
	})

	return router

}

// Run the server instance at given port
func (app *Application) Run() error {
	n := negroni.Classic()
	n.UseHandler(nosurf.New(app.Router()))
	n.Run(fmt.Sprintf(":%v", app.Options.Port))
	return nil
}

func indexPage(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	c := getRequestContext(ctx)

	var staticHost string

	if c.Options.Env == "dev" {
		staticHost = c.Options.DevServerURL
	}

	csrfToken := nosurf.Token(r)

	data := map[string]string{
		"staticHost": staticHost,
		"env":        c.Options.Env,
		"csrfToken":  csrfToken,
	}
	log.Println("indexpage", c.Options)
	c.Render.HTML(w, http.StatusOK, "index", data)
}

func getMovies(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := getRequestContext(ctx)
	movies, err := c.Store.GetAll()
	if err != nil {
		c.Err = err
		return
	}
	c.Render.JSON(w, http.StatusOK, movies)
}

func markSeen(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	c := getRequestContext(ctx)

	if err := c.Store.MarkSeen(pat.Param(ctx, "id")); err != nil {
		c.Err = err
		return
	}
	c.Render.Text(w, http.StatusOK, "Movie seen")
}

func getRandomMovie(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	c := getRequestContext(ctx)

	movie, err := c.Store.GetRandom()

	if err != nil {
		c.Err = err
		return
	}

	c.Render.JSON(w, http.StatusOK, movie)
}

func suggest(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	c := getRequestContext(ctx)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Err = err
		return
	}

	for {

		for {

			movie, err := c.Store.GetRandom()

			if err != nil {
				continue
			}

			if err := conn.WriteJSON(movie); err != nil {
				c.Err = err
				return
			}

			time.Sleep(socketWaitFor)
		}
	}

}

func getMovie(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := getRequestContext(ctx)
	movie, err := c.Store.Get(pat.Param(ctx, "id"))
	if err != nil {
		c.Err = err
		return
	}
	c.Render.JSON(w, http.StatusOK, movie)
}

func deleteMovie(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := getRequestContext(ctx)
	if err := c.Store.Delete(pat.Param(ctx, "id")); err != nil {
		c.Err = err
		return
	}
	c.Render.Text(w, http.StatusOK, "Movie deleted")
}

func addMovie(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	c := getRequestContext(ctx)

	d := &struct {
		Title string `valid:"required"`
	}{}
	if err := decode(r, d); err != nil {
		c.Err = err
		return
	}

	movie, err := c.OMDB.Find(d.Title)

	if err != nil {
		c.Err = err
		return
	}

	oldMovie, err := c.Store.Get(movie.ImdbID)

	if err == store.ErrMovieNotFound {

		if movie.Poster != "" {
			go func(url, imdbID string) {
				if filename, err := c.downloader.download(
					c.Options.StaticDir, url, imdbID); err != nil {
					log.Println(err)
					movie.Poster = ""
				} else {
					movie.Poster = filename
				}
				if err := c.Store.Save(movie); err != nil {
					log.Println(err)
				}
			}(movie.Poster, movie.ImdbID)
		}

		movie.Poster = "" // so we don't get a bad link

		if err := c.Store.Save(movie); err != nil {
			c.Err = err
			return
		}

		c.Render.JSON(w, http.StatusCreated, movie)
		return
	}

	if err != nil {
		c.Err = err
		return
	}

	c.Render.JSON(w, http.StatusOK, oldMovie)

}
