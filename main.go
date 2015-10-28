package main

import (
	"flag"
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

	// in a small app we could get away with globals
	// but here we'll use a global context object we can
	// inject into each request with all the useful things
	// we'll need

	// if we start using per-request context objects e.g.
	// logged in user, where it has to be threadsafe, then
	// use gorilla context.

	s := server.New(*env, db, staticURL, staticDir, devServerURL)

	chain := alice.New(nosurf.NewPure).Then(s.Router())
	http.ListenAndServe(":"+*port, chain)

}
