package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

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


// func (app *application) createPatientHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Name      string `json:"name"`
// 		Birthdate string `json:"birthdate"`
// 		Gender    string `json:"gender"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}

// 	patient := &model.Patient{
// 		Name:      input.Name,
// 		Birthdate: input.Birthdate,
// 		Gender:    input.Gender,
// 	}

// 	err = app.models.Patients.Insert(patient)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}

// 	app.respondWithJSON(w, http.StatusCreated, patient)
// }

// func (app *application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	param := vars["id"]

// 	id, err := strconv.Atoi(param)
// 	if err != nil || id < 1 {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
// 		return
// 	}

// 	patient, err := app.models.Patients.Get(id)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
// 		return
// 	}

// 	app.respondWithJSON(w, http.StatusOK, patient)
// }

// func (app *application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	param := vars["id"]

// 	id, err := strconv.Atoi(param)
// 	if err != nil || id < 1 {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
// 		return
// 	}

// 	patient, err := app.models.Patients.Get(id)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
// 		return
// 	}

// 	var input struct {
// 		Name      *string `json:"name"`
// 		Birthdate *string `json:"birthdate"`
// 		Gender    *string `json:"gender"`
// 	}

// 	err = app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}

// 	if input.Name != nil {
// 		patient.Name = *input.Name
// 	}

// 	if input.Birthdate != nil {
// 		patient.Birthdate = *input.Birthdate
// 	}

// 	if input.Gender != nil {
// 		patient.Gender = *input.Gender
// 	}

// 	err = app.models.Patients.Update(patient)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}

// 	app.respondWithJSON(w, http.StatusOK, patient)
// }

// func (app *application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	param := vars["id"]

// 	id, err := strconv.Atoi(param)
// 	if err != nil || id < 1 {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
// 		return
// 	}

// 	err = app.models.Patients.Delete(id)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}

// 	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
// }


func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
