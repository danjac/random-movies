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
	"net/http"
	"net/url"
)

type Movie struct {
	Title    string
	Actors   string
	Poster   string
	Year     string
	Plot     string
	Director string
	ImdbID   string `json:"imdbID"`
}

type MovieForm struct {
	Title string `valid:"required"`
}

func (f *MovieForm) Decode(r *http.Request) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(f); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return err
	}
	return nil
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

func main() {

	flag.Parse()

	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	render := render.New()

	// static content
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./dist/"))))

	// index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var staticHost string
		if *env == "dev" {
			staticHost = "http://localhost:8080"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, http.StatusOK, movie)
	}).Methods("GET")

	api.HandleFunc("/movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		movie, err := getMovie(db, mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.Text(w, http.StatusOK, "Deleted")
	}).Methods("DELETE")

	api.HandleFunc("/all/", func(w http.ResponseWriter, r *http.Request) {
		movies, err := getMovies(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if movie.Title != "" {
			if err := movie.Save(db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		render.JSON(w, http.StatusOK, movie)
	}).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":" + *port)

}
