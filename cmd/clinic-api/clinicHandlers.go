package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/validator"

	"github.com/gorilla/mux"
)

func (app *application) createClinicHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		City    string `json:"city"`
		Address string `json:"address"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	clinic := &model.Clinic{
		Name:    input.Name,
		City:    input.City,
		Address: input.Address,
	}

	err = app.models.Clinics.Insert(clinic)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"clinic": clinic}, nil)
}

func (app *application) searchClinicHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		City    string `json:"city"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.City = app.readStrings(qs, "city", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "city", "created_at", "updated_at",
		// descending sort values
		"-id", "-name", "-city", "-created_at", "-updated_at",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	clinics, metadata, err := app.models.Clinics.GetAll(input.Name, input.City, input.Filters)
	if err != nil {
		fmt.Println("We are in search clinic handler", "\nname: ", input.Name, "\ncity:", input.City, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"clinics": clinics, "metadata": metadata}, nil)
}

func (app *application) getClinicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid clinic ID")
		return
	}

	clinic, err := app.models.Clinics.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"clinic": clinic}, nil)
}

func (app *application) updateClinicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	fmt.Println(param)
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid clinic ID")
		fmt.Print("\nError: ", err)
		fmt.Print("\nid: ", id, "\n\n\n")
		return
	}

	clinic, err := app.models.Clinics.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name    *string `json:"name"`
		City    *string `json:"city"`
		Address *string `json:"address"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		clinic.Name = *input.Name
	}

	if input.City != nil {
		clinic.City = *input.City
	}

	if input.Address != nil {
		clinic.Address = *input.Address
	}

	err = app.models.Clinics.Update(clinic)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"clinic": clinic}, nil)
}

func (app *application) deleteClinicHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		fmt.Print("Error: ", err, ": ", id)
		return
	}

	err = app.models.Clinics.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
