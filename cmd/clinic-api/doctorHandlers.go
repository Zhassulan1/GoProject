package main

import "net/http"

func (app *application) createDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		name      string `json:"name"`
		specialty string `json:"specialty"`
	}
}
