package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Print("Received bad JSON")
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	log.Print("Successful JSON read")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PatientId string `json:"patientId"`
		DoctorId  string `json:"doctorId"`
		Date      string `json:"date"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Status    string `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	appointment := &model.Appointment{
		PatientId: input.PatientId,
		DoctorId:  input.DoctorId,
		Date:      input.Date,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    input.Status,
	}

	err = app.models.Appointments.Insert(appointment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, appointment)
}

func (app *application) getAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointment)
}

func (app *application) updateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		PatientId *string `json:"patientId"`
		DoctorId  *string `json:"doctorId"`
		Date      *string `json:"date"`
		StartTime *string `json:"startTime"`
		EndTime   *string `json:"endTime"`
		Status    *string `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.PatientId != nil {
		appointment.PatientId = *input.PatientId
	}

	if input.DoctorId != nil {
		appointment.DoctorId = *input.DoctorId
	}

	if input.Date != nil {
		appointment.Date = *input.Date
	}
	if input.StartTime != nil {
		appointment.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		appointment.EndTime = *input.EndTime
	}
	if input.Status != nil {
		appointment.Status = *input.Status
	}

	err = app.models.Appointments.Update(appointment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointment)
}

func (app *application) deleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	err = app.models.Appointments.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
