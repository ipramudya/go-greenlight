package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(rw http.ResponseWriter, r *http.Request, code int, message interface{}) {
	data := envelope{"error": message}

	err := app.writeJSON(rw, data, code, nil)
	if err != nil {
		app.logError(r, err)
		rw.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(rw, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(rw http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	app.errorResponse(rw, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(rw http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(rw, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequestResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(rw, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(rw http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(rw, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(rw http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(rw, r, http.StatusConflict, msg)
}

func (app *application) rateLimitExceedResponse(rw http.ResponseWriter, r *http.Request) {
	msg := "rate limit exceed"
	app.errorResponse(rw, r, http.StatusTooManyRequests, msg)
}
