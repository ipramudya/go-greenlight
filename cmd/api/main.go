package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/jsonlog"
	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenConn int
		maxIdleConn int
		maxIdleTime string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	setupFlag(&cfg)

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(logger, "", 0),
	}

	logger.PrintInfo("starting server...", map[string]string{
		"addr": s.Addr,
		"env":  app.env,
	})
	err = s.ListenAndServe()
	logger.PrintFatal(err, nil)
}
