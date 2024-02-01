package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const ApplicationJSON = "application/json"

func (app *application) readIDParam(r *http.Request) (int64, error) {
	/** When httprouter is parsing a request, any interpolated URL parameters
	 * 	stored in the request context. This method return a slice containing these parameter names and values
	 */
	params := httprouter.ParamsFromContext(r.Context())

	/** The value returned by ByName() is always a string. So we try to convert it
	 *	to a base 10 integer (with a bit size of 64), coz we know the ID can't be less than 1 and must int.
	 */
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(rw http.ResponseWriter, data envelope, code int, headers http.Header) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		rw.Header()[k] = v
	}

	json = append(json, '\n')
	rw.Header().Set("Content-Type", ApplicationJSON)
	rw.WriteHeader(code)
	rw.Write([]byte(json))

	return nil
}

func (app *application) readJSON(rw http.ResponseWriter, r *http.Request, destination interface{}) error {
	/* limit the size of the request body to 1MB */
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))

	/* This means that if the JSON from the client now includes any */
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(destination)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		/* There is a syntax problem with the JSON being decoded. */
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		/* There is a syntax problem with the JSON being decoded. */
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		/* A JSON value is not appropriate for the destination Go type. */
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at characted %d)", unmarshalTypeError.Offset)

		/* The JSON being decoded is empty */
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		/** Decode() will now return an error message in the format "json: unknown field "<name>"
		 * 	We check for this, extract the field name from the error & interpolate it to custom error message
		 */
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	/** Call Decode() again, using a pointer to an empty anonymous struct as the destination
	 *	If the request body only contained a single JSON value this will return an io.EOF error.
	 * 	So if we get anything else, we know that there is additional data in the request body
	 * 	and we return our own custom error message.
	 */
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
