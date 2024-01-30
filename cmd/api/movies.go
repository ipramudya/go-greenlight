package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ipramudya/go-greenlight/internal/data"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "create a new movie")
}

func (app *application) showMovieHandler(rw http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(rw, r)
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

	err = app.writeJSON(rw, movie, http.StatusOK, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
