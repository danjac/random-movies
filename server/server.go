package server

import (
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/errors"
	"github.com/danjac/random_movies/logger"
	"github.com/danjac/random_movies/models"
	"github.com/danjac/random_movies/utils"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"net/http"
)

func New(env string, db *database.DB, staticURL, staticDir, devServerURL string) *Server {
	return &Server{
		DB:           db,
		Env:          env,
		DevServerURL: devServerURL,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		Render:       render.New(),
		Log:          logger.New(),
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
	Log          *logger.Logger
}

func (s *Server) Abort(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(error).(type) {
	case errors.Error:
		s.Log.Error.Printf("HTTP %s %d: %s", e.Status(), e)
		http.Error(w, e.Error(), e.Status())
	default:
		s.Log.Error.Printf("%v: :%v", r, err)
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
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
	if err := s.DB.Del(mux.Vars(r)["id"]).Err(); err != nil {
		s.Abort(w, r, err)
		return
	}
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
	s.Log.Info.Printf("New movie %s added", movie)
	s.Render.JSON(w, http.StatusOK, movie)
}
