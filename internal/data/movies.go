package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ipramudya/go-greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`  // slices of genres for movie (romance, comedy, etc.)
	Version   int32     `json:"version"`           // The version number starts at 1 and will be incremented each
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

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id=$1`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, err
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5
	RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
