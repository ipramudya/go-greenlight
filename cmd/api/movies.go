package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ipramudya/go-greenlight/internal/data"
	"github.com/ipramudya/go-greenlight/internal/validator"
)

/** Endpont = "/v1/movies"
 *	Method = POST
 */
func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	/* decode destination */
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
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
		Runtime: data.Runtime(input.Runtime),
		Genres:  input.Genres,
	}

	/* field validation */
	vd := validator.New()

	if data.ValidateMovie(vd, movie); !vd.IsValid() {
		app.failedValidationResponse(rw, r, vd.Errors)
		return
	}

	fmt.Fprintf(rw, "%+v\n", input)
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

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"romance", "drama", "war"},
		Version:   1,
	}

	err = app.writeJSON(rw, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(rw, r, err)
	}
}
