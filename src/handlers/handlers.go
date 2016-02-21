package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"models"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/asaskevich/govalidator"

	"omdb"
	"store"

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

// New returns new AppConfig implementation
func New(db store.DB, options Options) *AppConfig {

	return &AppConfig{
		DB:      db,
		OMDB:    omdb.New(),
		Options: options,
		Render: render.New(render.Options{
			DisableHTTPErrorRendering: true,
		}),
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

// AppConfig is an instance of web app
type AppConfig struct {
	Options Options
	OMDB    omdb.Finder
	DB      store.DB
	Render  *render.Render
}

type handlerFunc func(*AppConfig, http.ResponseWriter, *http.Request) error

func (c *AppConfig) handler(fn handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(c, w, r); err != nil {
			if err == store.ErrMovieNotFound {
				http.NotFound(w, r)
				return
			}
			log.Println(err)
			http.Error(w, "An error has occurred", http.StatusInternalServerError)
		}
	})
}

// Router creates http handler
func (c *AppConfig) Router() http.Handler {

	router := mux.NewRouter()
	//router.StrictSlash(true)

	router.PathPrefix(c.Options.StaticURL).Handler(
		http.StripPrefix(c.Options.StaticURL,
			http.FileServer(http.Dir(c.Options.StaticDir))))

	api := router.PathPrefix("/api").Subrouter()

	api.Handle("/", c.handler(getRandomMovie)).Methods("GET")
	api.Handle("/", c.handler(addMovie)).Methods("POST")
	api.Handle("/suggest", c.handler(suggest)).Methods("GET")
	api.Handle("/movie/{id}", c.handler(getMovie)).Methods("GET")
	api.Handle("/movie/{id}", c.handler(deleteMovie)).Methods("DELETE")
	api.Handle("/seen/{id}", c.handler(markSeen)).Methods("PATCH")
	api.Handle("/all/", c.handler(getMovies)).Methods("GET")

	router.Handle("/{path:.*}", c.handler(indexPage)).Methods("GET")

	return router

}

// Run the server instance at given port
func (c *AppConfig) Run() error {
	n := negroni.Classic()
	n.UseHandler(nosurf.New(c.Router()))
	n.Run(fmt.Sprintf(":%v", c.Options.Port))
	return nil
}

func indexPage(c *AppConfig, w http.ResponseWriter, r *http.Request) error {

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
	return c.Render.HTML(w, http.StatusOK, "index", data)
}

func markSeen(c *AppConfig, w http.ResponseWriter, r *http.Request) error {
	if err := c.DB.MarkSeen(mux.Vars(r)["id"]); err != nil {
		return err
	}
	return c.Render.Text(w, http.StatusOK, "Movie seen")
}

func getRandomMovie(c *AppConfig, w http.ResponseWriter, r *http.Request) error {

	movie, err := c.DB.GetRandom()

	if err != nil {
		return err
	}

	return c.Render.JSON(w, http.StatusOK, movie)
}

func suggest(c *AppConfig, w http.ResponseWriter, r *http.Request) error {

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

func getMovie(c *AppConfig, w http.ResponseWriter, r *http.Request) error {
	movie, err := c.DB.Get(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	return c.Render.JSON(w, http.StatusOK, movie)
}

func deleteMovie(c *AppConfig, w http.ResponseWriter, r *http.Request) error {
	imdbID := mux.Vars(r)["id"]
	if err := c.DB.Delete(imdbID); err != nil {
		return err
	}
	return c.Render.Text(w, http.StatusOK, "Movie deleted")
}

func getMovies(c *AppConfig, w http.ResponseWriter, r *http.Request) error {
	movies, err := c.DB.GetAll()
	if err != nil {
		return err
	}
	return c.Render.JSON(w, http.StatusOK, movies)
}

func addMovie(c *AppConfig, w http.ResponseWriter, r *http.Request) error {
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

		if movie.Poster != "" {
			go func(movie *models.Movie) {
				if filename, err := downloadPoster(c.Options.StaticDir, movie.Poster, movie.ImdbID); err != nil {
					log.Println(err)
					movie.Poster = ""
				} else {
					movie.Poster = filename
				}
				if err := c.DB.Save(movie); err != nil {
					log.Println(err)
				}
			}(movie)
		}

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

func downloadPoster(staticDir, url, imdbID string) (string, error) {

	if url == "" {
		return "", nil
	}

	filename := fmt.Sprintf("%s.jpg", imdbID)

	imageDir := filepath.Join(staticDir, "images")
	if err := os.Mkdir(imageDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
	}

	imagePath := filepath.Join(imageDir, filename)

	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		return filename, err
	}

	// should probably check this
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filename, err
}
