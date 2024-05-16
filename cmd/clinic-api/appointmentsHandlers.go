package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/validator"
	"github.com/gorilla/mux"
)

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
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
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
		log.Print(err.Error())
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"appointment": appointment}, nil)
}

func (app *application) SearchAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PatientId string `json:"patientId"`
		DoctorId  string `json:"doctorId"`
		Date      string `json:"date"`
		Status    string `json:"status"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.PatientId = app.readStrings(qs, "patientId", "")
	input.DoctorId = app.readStrings(qs, "doctorId", "")
	input.Date = app.readStrings(qs, "adte", "")
	input.Status = app.readStrings(qs, "status", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "patientId", "doctorId", "date", "start_time", "end_time", "status", "created_at", "updated_at",
		// descending sort values
		"-id", "-patientId", "-doctorId", "-date", "-start_time", "-end_time", "-status", "-created_at", "-updated_at",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	appointments, metadata, err := app.models.Appointments.GetAll(input.PatientId, input.DoctorId, input.Date, input.Status, input.Filters)
	if err != nil {
		fmt.Println("We are in search appointment handler", "\npatientId: ", input.PatientId, "\ndoctorId:", input.DoctorId, "\ndate:", input.Date, "\nStatus:", input.Status, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"appointments": appointments, "metadata": metadata}, nil)
}

func (app *application) getAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"appointment": appointment}, nil)
}

func (app *application) updateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
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
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
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
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"appointment": appointment}, nil)
}

func (app *application) deleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Appointments.Delete(id)
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
