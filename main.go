package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/server"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"net/http"
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

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	db, err := database.New(database.DefaultConfig())

	if err != nil {
		log.Error("Bad Redis connection, shutting down...")
		panic(err)
	}

	s := server.New(db, log, &server.Config{
		Env:          *env,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		DevServerURL: devServerURL,
	})

	chain := alice.New(nosurf.NewPure).Then(s.Router())

	log.WithFields(logrus.Fields{
		"port": *port,
	}).Info("Server started")

	if err := http.ListenAndServe(":"+*port, chain); err != nil {
		panic(err)
	}

}
