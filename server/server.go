package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/decoders"
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/omdb"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
)

func New(db *database.DB, log *logrus.Logger, config *Config) *Server {
	return &Server{
		DB:     db,
		Render: render.New(),
		Log:    log,
		Config: config,
	}
}

// context globals (not threadsafe, so only store thread-safe objects here)
type Server struct {
	Config *Config
	Render *render.Render
	DB     *database.DB
	Log    *logrus.Logger
}

type Config struct {
	Env          string
	StaticURL    string
	StaticDir    string
	DevServerURL string
}

func (s *Server) Abort(w http.ResponseWriter, r *http.Request, err error) {
	logger := s.Log.WithFields(logrus.Fields{
		"URL":    r.URL,
		"Method": r.Method,
		"Error":  err,
	})
	var msg string
	switch e := err.(error).(type) {
	case errors.Error:
		msg = "HTTP Error"
		http.Error(w, e.Error(), e.Status())
	default:
		msg = "Internal Server Error"
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
	logger.Error(msg)
}

func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	// static content
	router.PathPrefix(
		s.Config.StaticURL).Handler(http.StripPrefix(
		s.Config.StaticURL, http.FileServer(http.Dir(s.Config.StaticDir))))

	// index page
	router.HandleFunc("/", s.indexPage).Methods("GET")

	// API calls
	api := router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/", s.getRandomMovie).Methods("GET")
	api.HandleFunc("/", s.addMovie).Methods("POST")
	api.HandleFunc("/movie/{id}", s.getMovie).Methods("GET")
	api.HandleFunc("/movie/{id}", s.deleteMovie).Methods("DELETE")
	api.HandleFunc("/all/", s.getMovies).Methods("GET")

	return router
}

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {

	var staticHost string

	if s.Config.Env == "dev" {
		staticHost = s.Config.DevServerURL
		s.Log.Info("Running development version")
	}

	csrfToken := nosurf.Token(r)

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        s.Config.Env,
		"csrfToken":  csrfToken,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}

func (s *Server) getRandomMovie(w http.ResponseWriter, r *http.Request) {

	movie, err := s.DB.GetRandomMovie()
	if err != nil {
		s.Abort(w, r, err)
		return
	}

	if movie == nil {
		s.Abort(w, r, errors.ErrHTTPNotFound)
		return
	}
	s.Render.JSON(w, http.StatusOK, movie)
}

func (s *Server) getMovie(w http.ResponseWriter, r *http.Request) {

	movie, err := s.DB.GetMovie(mux.Vars(r)["id"])
	if err != nil {
		s.Abort(w, r, err)
		return
	}
	if movie == nil {
		s.Abort(w, r, errors.ErrHTTPNotFound)
		return
	}
	s.Render.JSON(w, http.StatusOK, movie)
}

func (s *Server) deleteMovie(w http.ResponseWriter, r *http.Request) {
	imdbID := mux.Vars(r)["id"]
	if err := s.DB.Del(imdbID).Err(); err != nil {
		s.Abort(w, r, err)
		return
	}
	s.Log.WithFields(logrus.Fields{
		"imdbID": imdbID,
	}).Warn("Movie has been deleted")
	s.Render.Text(w, http.StatusOK, "Movie deleted")
}

func (s *Server) getMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := s.DB.GetMovies()
	if err != nil {
		s.Abort(w, r, err)
		return
	}
	s.Render.JSON(w, http.StatusOK, movies)
}

func (s *Server) addMovie(w http.ResponseWriter, r *http.Request) {

	f := &decoders.MovieDecoder{}
	if err := f.Decode(r); err != nil {
		s.Abort(w, r, errors.HTTPError{http.StatusBadRequest, err})
		return
	}

	movie, err := omdb.Search(f.Title)
	if err != nil {
		if err == omdb.ErrMovieNotFound {
			s.Log.WithFields(logrus.Fields{
				"title": f.Title,
			}).Warn("No movie found")
			s.Abort(w, r, errors.ErrHTTPNotFound)
		} else {
			s.Abort(w, r, err)
		}
		return
	}

	if err := s.DB.SaveMovie(movie); err != nil {
		s.Abort(w, r, err)
		return
	}
	s.Log.WithFields(logrus.Fields{
		"movie": movie,
	}).Info("New movie added")
	s.Render.JSON(w, http.StatusOK, movie)
}
