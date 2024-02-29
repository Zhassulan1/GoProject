package main

import (
	"log"
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
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	patient := &model.Patient{
		Name:      input.Name,
		Birthdate: input.Birthdate,
		Gender:    input.Gender,
	}

	err = app.models.Patients.Insert(patient)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, patient)
}

func (app *application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := app.models.Patients.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patient)
}

func (app *application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Patient Update Started")

	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := app.models.Patients.Get(id)
	if err != nil {
		log.Print("Could Not Upadate Patient")
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name      *string `json:"name"`
		Birthdate *string `json:"birthdate"`
		Gender    *string `json:"gender"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patient)
}

func (app *application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	err = app.models.Patients.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
