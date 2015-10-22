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
	"runtime"
)

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

func (l *Logger) WriteErr(w http.ResponseWriter, err error) {
	_, fn, line, _ := runtime.Caller(1)
	l.Error.Printf("%s:%d:%v", fn, line, err)
	http.Error(w, "Sorry, an error has occurred", http.StatusInternalServerError)
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

func (m *Movie) Save(db *redis.Client) error {
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

type Env struct {
	*render.Render
	DB  *redis.Client
	Log *Logger
}

type HandlerFunc func(e *Env, w http.ResponseWriter, r *http.Request) error

type Handler struct {
	*Env
	H HandlerFunc
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			h.Env.Log.Error.Printf("HTTP %d: %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			_, fn, line, _ := runtime.Caller(1)
			h.Env.Log.Error.Printf("%s:%d:%v", fn, line, err)
			http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
		}
	}
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

	return movie, nil
}

func getRandomMovie(db *redis.Client) (*Movie, error) {
	imdbID, err := db.RandomKey().Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return getMovie(db, imdbID)
}

func getMovie(db *redis.Client, imdbID string) (*Movie, error) {
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

func getMovies(db *redis.Client) ([]*Movie, error) {
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

	db := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	render := render.New()
	logger := newLogger()

	// static content
	router.PathPrefix(
		staticURL).Handler(http.StripPrefix(
		staticURL, http.FileServer(http.Dir(staticDir))))

	// index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var staticHost string
		if *env == "dev" {
			staticHost = devServerURL
		}

		ctx := map[string]string{
			"staticHost": staticHost,
			"env":        *env,
		}

		render.HTML(w, http.StatusOK, "index", ctx)
	})

	// API calls
	api := router.PathPrefix("/api/").Subrouter()

	e := &Env{render, db, logger}

	api.Handle("/", Handler{e, func(e *Env, w http.ResponseWriter, r *http.Request) error {
		movie, err := getRandomMovie(e.DB)
		if err != nil {
			return err
		}
		if movie == nil {
			return errHTTPNotFound
		}
		e.JSON(w, http.StatusOK, movie)
		return nil
	}}).Methods("GET")

	api.Handle("/movie/{id}", Handler{e, func(e *Env, w http.ResponseWriter, r *http.Request) error {
		movie, err := getMovie(e.DB, mux.Vars(r)["id"])
		if err != nil {
			return err
		}
		if movie == nil {
			return errHTTPNotFound
		}
		e.JSON(w, http.StatusOK, movie)
		return nil
	}}).Methods("GET")

	api.Handle("/movie/{id}", Handler{e, func(e *Env, w http.ResponseWriter, r *http.Request) error {
		if err := db.Del(mux.Vars(r)["id"]).Err(); err != nil {
			return err
		}
		e.Text(w, http.StatusOK, "Movie deleted")
		return nil
	}}).Methods("DELETE")

	api.Handle("/all/", Handler{e, func(e *Env, w http.ResponseWriter, r *http.Request) error {
		movies, err := getMovies(e.DB)
		if err != nil {
			return err
		}
		e.JSON(w, http.StatusOK, movies)
		return nil
	}}).Methods("GET")

	api.Handle("/", Handler{e, func(e *Env, w http.ResponseWriter, r *http.Request) error {

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

		if err := movie.Save(e.DB); err != nil {
			return err
		}
		logger.Info.Printf("New movie %s added", movie)
		e.JSON(w, http.StatusOK, movie)
		return nil

	}}).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + *port)

}
