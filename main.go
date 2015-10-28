package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/redis.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

func (log *Logger) Handle(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(error).(type) {
	case Error:
		log.Error.Printf("HTTP %d: %s", e.Status(), e)
		http.Error(w, e.Error(), e.Status())
	default:
		_, fn, line, _ := runtime.Caller(1)
		log.Error.Printf("%s:%d:%v", fn, line, err)
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
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

type AppContext struct {
	Env    string
	Render *render.Render
	DB     *DB
	Log    *Logger
}

func getAppContext(r *http.Request) *AppContext {
	return context.Get(r, "appContext").(*AppContext)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)

	var staticHost string
	if c.Env == "dev" {
		staticHost = devServerURL
	}

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        c.Env,
	}

	c.Render.HTML(w, http.StatusOK, "index", ctx)
}

func getRandomMovie(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)
	movie, err := c.DB.GetRandomMovie()
	if err != nil {
		c.Log.Handle(w, r, err)
		return
	}
	if movie == nil {
		c.Log.Handle(w, r, errHTTPNotFound)
		return
	}
	c.Render.JSON(w, http.StatusOK, movie)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)
	movie, err := c.DB.GetMovie(mux.Vars(r)["id"])
	if err != nil {
		c.Log.Handle(w, r, err)
		return
	}
	if movie == nil {
		c.Log.Handle(w, r, errHTTPNotFound)
		return
	}
	c.Render.JSON(w, http.StatusOK, movie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)
	if err := c.DB.Del(mux.Vars(r)["id"]).Err(); err != nil {
		c.Log.Handle(w, r, err)
		return
	}
	c.Render.Text(w, http.StatusOK, "Movie deleted")
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)
	movies, err := c.DB.GetMovies()
	if err != nil {
		c.Log.Handle(w, r, err)
		return
	}
	c.Render.JSON(w, http.StatusOK, movies)
}

func addMovie(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)
	f := &MovieForm{}
	if err := f.Decode(r); err != nil {
		c.Log.Handle(w, r, HTTPError{http.StatusBadRequest, err})
		return
	}

	movie, err := getMovieFromOMDB(f.Title)
	if err != nil {
		c.Log.Handle(w, r, err)
		return
	}

	if movie.ImdbID == "" {
		c.Log.Handle(w, r, errHTTPNotFound)
		return
	}

	if err := movie.Save(c.DB); err != nil {
		c.Log.Handle(w, r, err)
		return
	}
	c.Log.Info.Printf("New movie %s added", movie)
	c.Render.JSON(w, http.StatusOK, movie)
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

	c := &AppContext{
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
	router.HandleFunc("/", indexPage).Methods("GET")

	// API calls
	api := router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/", getRandomMovie).Methods("GET")
	api.HandleFunc("/", addMovie).Methods("POST")
	api.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	api.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
	api.HandleFunc("/all/", getMovies).Methods("GET")

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, "appContext", c)
		next(w, r)
	}))
	n.UseHandler(router)
	n.Run(":" + *port)

}
