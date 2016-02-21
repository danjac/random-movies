package main

import (
	"flag"

	"handlers"
	"store"
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

	db, err := store.New(store.DefaultConfig())

	if err != nil {
		panic(err)
	}

	cfg := handlers.New(db, handlers.Options{
		Env:          *env,
		StaticURL:    staticURL,
		StaticDir:    staticDir,
		DevServerURL: devServerURL,
		Port:         6060,
	})

	if err := cfg.Run(); err != nil {
		panic(err)
	}

}
