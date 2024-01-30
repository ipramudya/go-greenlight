package main

import (
	"net/http"
)

func (app *application) healthHandler(rw http.ResponseWriter, r *http.Request) {
	data := envelope{
		"health": envelope{
			"status": "available",
			"system_info": map[string]string{
				"environtment": app.config.env,
				"version":      Version,
			},
		},
	}

	err := app.writeJSON(rw, data, http.StatusOK, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
