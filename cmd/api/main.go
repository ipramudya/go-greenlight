package main

import (
	"os"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/jsonlog"
	"github.com/ipramudya/go-greenlight/internal/mailer"
	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type application struct {
	config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
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
		mailer: *mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}
