//Filename: cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// create a method that returns a http router
func (app *application) routes() *httprouter.Router {
	//Create a new httrouter router instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedesponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/list", app.createListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/list/:id", app.showListHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/list/:id", app.updateListHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/list/:id", app.deleteListHandler)

	return router
}
