package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// when we receive a SIGINT or SIGTERM signal, we instruct our server to stop accepting
	// any new HTTP requests, and give any in-flight requests a ‘grace
	// period’ of 5 seconds to complete before the application is terminated.
	go func() {
		quit := make(chan os.Signal, 1)

		// catch two signal and store that into channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// read any signal, this code wil block below its until a signal received
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown will return nil if the graceful shutdown was successful, or
		// an error (happen bcs of a problem closing listeners, or didn't complete before 5 sec context deadline is hit)
		shutdownError <- server.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  app.env,
	})

	// calling shutdown() on our server (which is via goroutine above) will cause ListenAndServe() to immediately return a http.ErrServerClosed.
	// then if error caught, it's indicated that graceful shutdown has started.
	// so we check  only returning the error IF IT IS NOT http.ErrServerClosed, then we handled it by returning that other error.
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// we wait to receive the return value from Shutdown()
	// (as if we know the error is http.ErrServerClosed) on the shutdownError channel.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": server.Addr,
	})

	return nil
}
