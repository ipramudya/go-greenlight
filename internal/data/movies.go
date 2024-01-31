package data

import (
	"time"

	"github.com/ipramudya/go-greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`         // slices of genres for movie (romance, comedy, etc.)
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each
}

func ValidateMovie(vd *validator.Validator, movie *Movie) {
	/* title validation */
	vd.Check(movie.Title != "", "title", "must be provided")
	vd.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	/* year validation */
	vd.Check(movie.Year != 0, "year", "must be provided")
	vd.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	vd.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	/* runtime validation */
	vd.Check(movie.Runtime != 0, "runtime", "must be provided")
	vd.Check(movie.Runtime > 0, "runtime", "must be a positive")
	/* genres validation */
	vd.Check(movie.Genres != nil, "genres", "must be provided")
	vd.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	vd.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	/* unique validation */
	vd.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicated genre")
}
