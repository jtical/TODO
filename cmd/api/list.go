//Filename: cmd/api/list.go

package main

import (
	"fmt"
	"net/http"
)

// createListHandler for the "POST /v1/list" endpoint
func (app *application) createListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "creating a new todo list...")
}

// showListHandler for the "GET /v1/list" endpoint
func (app *application) showListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	//display the list id
	fmt.Fprintf(w, "show the details for list %d\n", id)
}
