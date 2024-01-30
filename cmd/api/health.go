package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "status: available")
	fmt.Fprintf(rw, "environtment: %s\n", app.config.env)
	fmt.Fprintf(rw, "version: %s\n", Version)
}
