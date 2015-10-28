package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// context globals (not threadsafe, so only store thread-safe objects here)
type AppConfig struct {
	Env    string
	Render *render.Render
	DB     *DB
	Log    *Logger
}

// per-request context; we'll stick threadsafe stuff in here at some point
type Context struct {
	*AppConfig
}

func (app *AppConfig) Handler(h HandlerFunc) http.Handler {
	return Handler{app, h}
}

type HandlerFunc func(c *Context, w http.ResponseWriter, r *http.Request) error

type Handler struct {
	*AppConfig
	H HandlerFunc
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{h.AppConfig}
	err := h.H(c, w, r)
	if err != nil {
		switch e := err.(error).(type) {
		case Error:
			h.Log.Error.Printf("HTTP %s %d: %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			h.Log.Error.Printf("%v: :%v", r, err)
			http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
		}
	}

}

type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

func newLogger() *Logger {

	flag := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		log.New(os.Stdout, "DEBUG: ", flag),
		log.New(os.Stdout, "INFO: ", flag),
		log.New(os.Stdout, "WARN: ", flag),
		log.New(os.Stderr, "ERROR: ", flag),
	}
}

type Movie struct {
	Title    string
	Actors   string
	Poster   string
	Year     string
	Plot     string
	Director string
	Rating   string `json:"imdbRating"`
	ImdbID   string `json:"imdbID"`
}

func (m *Movie) String() string {
	return m.Title
}

func (m *Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Movie) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Movie) Save(db *DB) error {
	return db.Set(m.ImdbID, m, 0).Err()
}

type MovieForm struct {
	Title string `valid:"required"`
}

func (f *MovieForm) Decode(r *http.Request) error {
	return decode(r, f)
}

type Error interface {
	error
	Status() int
}

type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func (e HTTPError) Status() int {
	return e.Code
}

var errHTTPNotFound = HTTPError{http.StatusNotFound, errors.New("Not found")}

func indexPage(c *Context, w http.ResponseWriter, r *http.Request) error {

	var staticHost string

	if c.Env == "dev" {
		staticHost = devServerURL
	}

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        c.Env,
	}

	c.Render.HTML(w, http.StatusOK, "index", ctx)
	return nil
}

func getRandomMovie(c *Context, w http.ResponseWriter, r *http.Request) error {
	movie, err := c.DB.GetRandomMovie()
	if err != nil {
		return err
	}
	if movie == nil {
		return errHTTPNotFound
	}
	c.Render.JSON(w, http.StatusOK, movie)
	return nil
}

func getMovie(c *Context, w http.ResponseWriter, r *http.Request) error {
	movie, err := c.DB.GetMovie(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	if movie == nil {
		return errHTTPNotFound
	}
	c.Render.JSON(w, http.StatusOK, movie)
	return nil
}

func deleteMovie(c *Context, w http.ResponseWriter, r *http.Request) error {
	if err := c.DB.Del(mux.Vars(r)["id"]).Err(); err != nil {
		return err
	}
	c.Render.Text(w, http.StatusOK, "Movie deleted")
	return nil
}

func getMovies(c *Context, w http.ResponseWriter, r *http.Request) error {
	movies, err := c.DB.GetMovies()
	if err != nil {
		return err
	}
	c.Render.JSON(w, http.StatusOK, movies)
	return nil
}

func addMovie(c *Context, w http.ResponseWriter, r *http.Request) error {
	f := &MovieForm{}
	if err := f.Decode(r); err != nil {
		return HTTPError{http.StatusBadRequest, err}
	}

	movie, err := getMovieFromOMDB(f.Title)
	if err != nil {
		return err
	}

	if movie.ImdbID == "" {
		return errHTTPNotFound
	}

	if err := movie.Save(c.DB); err != nil {
		return err
	}
	c.Log.Info.Printf("New movie %s added", movie)
	c.Render.JSON(w, http.StatusOK, movie)
	return nil
}

// decodes JSON body of request and runs through validator
func decode(r *http.Request, data interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return err
	}
	return nil
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

	if movie.ImdbID == "" {
		return nil, errors.New("Movie not found")
	}

	return movie, nil
}

type DB struct {
	*redis.Client
}

func (db *DB) GetRandomMovie() (*Movie, error) {
	imdbID, err := db.RandomKey().Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return db.GetMovie(imdbID)
}

func (db *DB) GetMovie(imdbID string) (*Movie, error) {
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

func (db *DB) GetMovies() ([]*Movie, error) {
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

const (
	staticURL    = "/static/"
	staticDir    = "./dist/"
	devServerURL = "http://localhost:8080"
	redisAddr    = "localhost:6379"
)

func main() {

	flag.Parse()

	db := &DB{redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})}

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	// in a small app we could get away with globals
	// but here we'll use a global context object we can
	// inject into each request with all the useful things
	// we'll need

	// if we start using per-request context objects e.g.
	// logged in user, where it has to be threadsafe, then
	// use gorilla context.

	app := &AppConfig{
		Env:    *env,
		DB:     db,
		Render: render.New(),
		Log:    newLogger(),
	}

	router := mux.NewRouter()

	// static content
	router.PathPrefix(
		staticURL).Handler(http.StripPrefix(
		staticURL, http.FileServer(http.Dir(staticDir))))

	// index page
	router.Handle("/", app.Handler(indexPage)).Methods("GET")

	// API calls
	api := router.PathPrefix("/api/").Subrouter()

	api.Handle("/", app.Handler(getRandomMovie)).Methods("GET")
	api.Handle("/", app.Handler(addMovie)).Methods("POST")
	api.Handle("/movie/{id}", app.Handler(getMovie)).Methods("GET")
	api.Handle("/movie/{id}", app.Handler(deleteMovie)).Methods("DELETE")
	api.Handle("/all/", app.Handler(getMovies)).Methods("GET")

	n := negroni.Classic()

	n.UseHandler(router)
	n.Run(":" + *port)

}
