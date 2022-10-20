//Filename: cmd/api/list.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"todo.joelical.net/internal/data"
	"todo.joelical.net/internal/validator"
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
	//copy the values from the input struct to a new lists struct
	list := &data.List{
		Name:   input.Name,
		Task:   input.Task,
		Status: input.Status,
	}

	//Initialize a new validator instance
	v := validator.New()

	//check the map to determain if there were any validation errors
	if data.ValidateList(v, list); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	//create a list
	err = app.models.List.Insert(list)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//create a location header for the newly created resource
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/list/%d", list.ID))
	//write the response with 201 -created status code with the body being the list data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"list": list}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
