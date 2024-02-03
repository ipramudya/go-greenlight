package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		// quit := make(chan os.Signal, 1)
		quit := make(chan os.Signal, 1)

		// catch two signal and store that into channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read any signal, this code wil block below its until a signal received
		s := <-quit

		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		// exit app with 0 (success) status code
		os.Exit(0)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": s.Addr,
		"env":  app.env,
	})

	return s.ListenAndServe()
}
