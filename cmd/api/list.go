//Filename: cmd/api/list.go

package main

import (
	"errors"
	"fmt"
	"net/http"

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

	//Fetch the specific list
	list, err := app.models.List.Get(id)
	//handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//write the data returned by get()
	err = app.writeJSON(w, http.StatusOK, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// updateListHandler for the "PUT /v1/list/:id" endpoint
func (app *application) updateListHandler(w http.ResponseWriter, r *http.Request) {
	//this method does a complete replacement
	//get the id for the list that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//fetch the orginal record from the database
	list, err := app.models.List.Get(id)
	//handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//create an input struct to hold data read in from the user
	// our target decode destination
	var input struct {
		Name   string `json:"name"`
		Task   string `json:"task"`
		Status string `json:"status"`
	}
	//initialize a new json.decode instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//copy and update the values in the list variable using the fields in the input struct
	list.Name = input.Name
	list.Task = input.Task
	list.Status = input.Status
	//perform validation on the updated list. if validation fails, then we send a 422 - unprocessable entity response to the user
	//Initialize a new validator instance
	v := validator.New()

	//check the map to determain if there were any validation errors
	if data.ValidateList(v, list); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	//pass the updated list record to the update() method
	err = app.models.List.Update(list)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//write the data returned by get()
	err = app.writeJSON(w, http.StatusOK, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// deleteListHandler for the "DELETE /v1/list/:id" endpoint
func (app *application) deleteListHandler(w http.ResponseWriter, r *http.Request) {
	//gets the id for the list that will be deleted
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//delete the list from the database. sends a 404 not found status code to the user if there is no matching record.
	err = app.models.List.Delete(id)
	//handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//return a 200 status ok to the user with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "list successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
