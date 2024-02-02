package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				app.serverErrorResponse(rw, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	limitter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !limitter.Allow() {
			app.rateLimitExceedResponse(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
