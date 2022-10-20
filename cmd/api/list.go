//Filename: cmd/api/list.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"todo.joelical.net/internal/data"
)

// createListHandler for the "POST /v1/list" endpoint
func (app *application) createListHandler(w http.ResponseWriter, r *http.Request) {
	// our target decode destination
	var input struct {
		Name   string `json:"name"`
		Task   string `json:"task"`
		Status string `json:"status"`
	}
	//initialize a new json.decode instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//display the request
	fmt.Fprintf(w, "%+v\n", input)
}

// showListHandler for the "GET /v1/list" endpoint
func (app *application) showListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//create a new instance of the List struct containing the ID we extracted from our url and some sample data
	list := data.List{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Study",
		Task:      "study for algebra test",
		Status:    "Completed",
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
