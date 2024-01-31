package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const Version = "1.0.0"

type config struct {
	port int    // network port address
	env  string // (development, stagging, production, etc)
}

type application struct {
	config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environtment (development|stagging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on port %s", cfg.env, s.Addr)
	err := s.ListenAndServe()
	logger.Fatal(err)
}
