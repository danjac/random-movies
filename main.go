package main

import (
	"encoding/json"
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

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		movie, err := getRandomMovie(db)
		if err != nil {
			logger.WriteErr(w, err)
			return
		}
		if movie == nil {
			http.NotFound(w, r)
			return
		}
		render.JSON(w, http.StatusOK, movie)
	}).Methods("GET")

	api.HandleFunc("/movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		movie, err := getMovie(db, mux.Vars(r)["id"])
		if err != nil {
			logger.WriteErr(w, err)
			return
		}
		if movie == nil {
			http.NotFound(w, r)
			return
		}
		render.JSON(w, http.StatusOK, movie)
	}).Methods("GET")

	api.HandleFunc("/movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Del(mux.Vars(r)["id"]).Err(); err != nil {
			logger.WriteErr(w, err)
			return
		}
		render.Text(w, http.StatusOK, "Movie deleted")
	}).Methods("DELETE")

	api.HandleFunc("/all/", func(w http.ResponseWriter, r *http.Request) {
		movies, err := getMovies(db)
		if err != nil {
			logger.WriteErr(w, err)
			return
		}
		render.JSON(w, http.StatusOK, movies)
	}).Methods("GET")

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		f := &MovieForm{}
		if err := f.Decode(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		movie, err := getMovieFromOMDB(f.Title)
		if err != nil {
			logger.WriteErr(w, err)
			return
		}

		if movie.ImdbID == "" {
			logger.Warn.Printf("No movie found for title %s", f.Title)
			http.NotFound(w, r)
			return
		}

		if err := movie.Save(db); err != nil {
			logger.WriteErr(w, err)
			return
		}
		logger.Info.Printf("New movie %s added", movie)
		render.JSON(w, http.StatusOK, movie)

	}).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + *port)

}
