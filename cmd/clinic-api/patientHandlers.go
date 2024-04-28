package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"

	"github.com/gorilla/mux"
)

func (app *application) createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Birthdate string `json:"birthdate"`
		Gender    string `json:"gender"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	patient := &model.Patient{
		Name:      input.Name,
		Birthdate: input.Birthdate,
		Gender:    input.Gender,
	}

	err = app.models.Patients.Insert(patient)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"patient": patient}, nil)
}

func (app *application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := app.models.Patients.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"patient": patient}, nil)
}

func (app *application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := app.models.Patients.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name      *string `json:"name"`
		Birthdate *string `json:"birthdate"`
		Gender    *string `json:"gender"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		patient.Name = *input.Name
	}

	if input.Birthdate != nil {
		patient.Birthdate = *input.Birthdate
	}

	if input.Gender != nil {
		patient.Gender = *input.Gender
	}

	err = app.models.Patients.Update(patient)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"patient": patient}, nil)
}

func (app *application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Patients.Delete(id)
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
