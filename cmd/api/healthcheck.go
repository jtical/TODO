//Filename: cmd/api/healthcheck.go

package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//create a map to hold our healthcheck data
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	//convert our map into a json object
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The Server encountered a prolem and could not Process request", http.StatusInternalServerError)
		return
	}
	//add a newline to make view on terminaL better
	js = append(js, '\n')
	//specify that we will serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	//write the [] byte slice containing the JSON response body
	w.Write(js)
}
