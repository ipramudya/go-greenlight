package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/validator"
)

/** Endpont = "/v1/movies"
 *	Method = POST
 */
func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	/* decode destination */
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	/** Importantly, notice that when we call Decode() we pass a *pointer* to the input
	 *	struct as the target decode destination. This must non-nil pointer as decoded target
	 */
	err := app.readJSON(rw, r, &input)
	if err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	/* field validation */
	vd := validator.New()

	if data.ValidateMovie(vd, movie); !vd.IsValid() {
		app.failedValidationResponse(rw, r, vd.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(rw, envelope{"movie": movie}, http.StatusCreated, headers)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}

/** Endpont = "/v1/movies/:id"
 *	Method = GET
 */
func (app *application) showMovieHandler(rw http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	err = app.writeJSON(rw, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(rw, r, err)
	}
}

/** Endpont = "/v1/movies/:id"
 *	Method = PUT
 */
func (app *application) updateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(rw, r, &input)

	if err != nil {
		app.badRequestResponse(rw, r, err)
		return
	}
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
		return
	}

	err = app.writeJSON(rw, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}

/** Endpont = "/v1/movies/:id"
 *	Method = DELETE
 */
func (app *application) deleteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(rw, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(rw, r)
		default:
			app.serverErrorResponse(rw, r, err)
		}
		return
	}

	err = app.writeJSON(rw, envelope{"message": "movie successfully deleted"}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}
