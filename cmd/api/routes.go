package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	r := httprouter.New()

	/* error handlers */
	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	/* routes handlers */
	r.HandlerFunc(http.MethodGet, "/v1/health", app.healthHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	r.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)
	r.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	return r
}
