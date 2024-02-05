package main

import (
	"errors"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	app.background(
		func() {
			data := map[string]interface{}{
				"activationToken": token.Plaintext,
				"userID":          token.UserID,
			}
			err := app.mailer.Send(user.Email, "user_welcome.tmpl", data)
			if err != nil {
				app.logger.PrintError(err, nil)
				return
			}
		},
	)

	err = app.writeJSON(rw, envelope{"user": user}, http.StatusAccepted, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}

func (app *application) activateUserHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := app.readJSON(rw, r, &input); err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)
		return
	}

	user, err := app.models.Tokens.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(rw, r, v.Errors)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	err = app.writeJSON(rw, envelope{"user": user}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}

}
