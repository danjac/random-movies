package main

import (
	"flag"

	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/server"
	"github.com/justinas/nosurf"
)

var (
	env  = flag.String("env", "prod", "environment ('prod' or 'dev')")
	port = flag.String("port", "4000", "server port")
)

const (
	staticURL    = "/static/"
	staticDir    = "./dist/"
	devServerURL = "http://localhost:8080"
)

func main() {

	flag.Parse()

	db, err := database.New(database.DefaultConfig())

	if err != nil {
		panic(err)
	}

	s := server.New(db, &server.Config{
		Env:          *env,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		DevServerURL: devServerURL,
		Port:         6060,
	})

	server := s.Router()
	server.Handler = nosurf.NewPure(server.Handler)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
