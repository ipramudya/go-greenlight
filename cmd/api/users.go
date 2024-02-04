package main

import (
	"errors"
	"net/http"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/validator"
)

func (app *application) registerUserHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(rw, r, &input)
	if err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)
		return
	}

	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exist")
			app.failedValidationResponse(rw, r, v.Errors)
		default:
			app.serverErrorResponse(rw, r, err)
		}

		return
	}

	err = app.writeJSON(rw, envelope{"user": user}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}
