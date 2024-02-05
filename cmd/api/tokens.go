package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(rw, r, &input)
	if err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	v := validator.New()

	if !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(rw, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	err = app.writeJSON(rw, envelope{"authentication_token": token}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}
