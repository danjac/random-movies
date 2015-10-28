package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/danjac/random_movies/database"
	"github.com/danjac/random_movies/server"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"gopkg.in/redis.v3"
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
	redisAddr    = "localhost:6379"
)

func main() {

	flag.Parse()

	db := &database.DB{redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})}

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	s := server.New(*env, db, log, staticURL, staticDir, devServerURL)

	chain := alice.New(nosurf.NewPure).Then(s.Router())
	http.ListenAndServe(":"+*port, chain)

}
