package main

import (
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"

	"github.com/gorilla/mux"
)


func (app *application) createDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string `json:"name"`
		Specialty  string `json:"specialty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	doctor := &model.Doctor{
		Name:      input.Name,
		Specialty: input.Specialty,
	}

	err = app.models.Doctors.Insert(doctor)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, doctor)
}

func (app *application) getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) updateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name      *string `json:"name"`
		Specialty *string `json:"specialty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		doctor.Name = *input.Name
	}

	if input.Specialty != nil {
		doctor.Specialty = *input.Specialty
	}

	
	err = app.models.Doctors.Update(doctor)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) deleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	err = app.models.Doctors.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}