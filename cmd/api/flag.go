package main

import (
	"flag"
	"os"
)

func setupFlag(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environtment (development|stagging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-cons", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConn, "db-max-idle-cons", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "122595c8bcd3c9", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "aee8a44eed97d6", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "ipramudya.dev@gmail.com", "SMTP sender")
	flag.Parse()
}
