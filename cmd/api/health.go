package main

import (
	"encoding/json"
	"net/http"
)

const ApplicationJSON = "application/json"

func (app *application) healthHandler(rw http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":       "available",
		"environtment": app.config.env,
		"version":      Version,
	}

	json, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	json = append(json, '\n')
	rw.Header().Set("Content-Type", ApplicationJSON)
	rw.Write([]byte(json))
}
