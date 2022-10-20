// Filename - cmd/api/errors.go
package main

import (
	"fmt"
	"net/http"
)

// create method to log errors
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// we want to display json formated error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//created json resposne
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// server error response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log the error
	app.logError(r, err)
	//prepare the response with error
	message := "the server encounter a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// not found response
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	//create error message
	message := "the requested page could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// a method not allowed response
func (app *application) methodNotAllowedesponse(w http.ResponseWriter, r *http.Request) {
	//create error message
	message := fmt.Sprintf("the %s method is not allowed for this resource", r.Method) // use formated string to allow message to be dynamic
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// a method for bad request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	//create our message
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
