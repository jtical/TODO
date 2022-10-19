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
	fmt.Fprintln(w, "creating a new todo list...")
}

// showListHandler for the "GET /v1/list" endpoint
func (app *application) showListHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
	err = app.writeJSON(w, http.StatusOK, list, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The serever encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
