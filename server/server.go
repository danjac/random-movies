package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/models"
	"github.com/danjac/random_movies/utils"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
)

func New(env string, db *database.DB, log *logrus.Logger, staticURL, staticDir, devServerURL string) *Server {
	return &Server{
		DB:           db,
		Env:          env,
		DevServerURL: devServerURL,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		Render:       render.New(),
		Log:          log,
	}
}

// context globals (not threadsafe, so only store thread-safe objects here)
type Server struct {
	Env          string
	StaticURL    string
	StaticDir    string
	DevServerURL string
	Render       *render.Render
	DB           *database.DB
	Log          *logrus.Logger
}

func (s *Server) Abort(w http.ResponseWriter, r *http.Request, err error) {
	logger := s.Log.WithFields(logrus.Fields{
		"Request": r,
		"Error":   err,
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
		s.StaticURL).Handler(http.StripPrefix(
		s.StaticURL, http.FileServer(http.Dir(s.StaticDir))))

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

	if s.Env == "dev" {
		staticHost = s.DevServerURL
	}

	csrfToken := nosurf.Token(r)

	ctx := map[string]string{
		"staticHost": staticHost,
		"env":        s.Env,
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

	f := &models.MovieForm{}
	if err := f.Decode(r); err != nil {
		s.Abort(w, r, errors.HTTPError{http.StatusBadRequest, err})
		return
	}

	movie, err := utils.GetMovieFromOMDB(f.Title)
	if err != nil {
		s.Abort(w, r, err)
		return
	}

	if movie.ImdbID == "" {
		s.Abort(w, r, errors.ErrHTTPNotFound)
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
