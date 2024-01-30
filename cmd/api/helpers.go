package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
