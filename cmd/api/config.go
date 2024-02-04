package main

type config struct {
	port    int
	env     string
	db      db
	limiter limiter
	smtp    smtp
}

type db struct {
	dsn         string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

type limiter struct {
	rps     float64
	burst   int
	enabled bool
}

type smtp struct {
	host     string
	port     int
	username string
	password string
	sender   string
}
