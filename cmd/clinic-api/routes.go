package main

import (
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// CLinic Singleton
	// Create a new appointment
	v1.HandleFunc("/appointments", app.createAppointmentHandler).Methods("POST")
	// Get a doctors list by pagination and filters
	v1.HandleFunc("/appointments", app.requirePermissions("appointments:read", app.SearchAppointmentHandler)).Methods("GET")
	// Get a specific appointment
	v1.HandleFunc("/appointments/{id:[0-9]+}", app.requirePermissions("appointments:read", app.getAppointmentHandler)).Methods("GET")
	// Update a specific appointment
	v1.HandleFunc("/appointments/{id:[0-9]+}", app.requirePermissions("appointments:write", app.updateAppointmentHandler)).Methods("PUT")
	// Delete a specific appointment
	v1.HandleFunc("/appointments/{id:[0-9]+}", app.requirePermissions("appointments:write", app.deleteAppointmentHandler)).Methods("DELETE")

	// Create a new doctor
	v1.HandleFunc("/doctors", app.createDoctorHandler).Methods("POST")
	// Get a doctors list by pagination and filters
	v1.HandleFunc("/doctors", app.requirePermissions("doctors:read", app.SearchDoctorHandler)).Methods("GET")
	// Get a specific doctor
	v1.HandleFunc("/doctors/{id:[0-9]+}", app.requirePermissions("doctors:read", app.getDoctorHandler)).Methods("GET")
	// Update a specific doctor
	v1.HandleFunc("/doctors/{id:[0-9]+}", app.requirePermissions("doctors:write", app.updateDoctorHandler)).Methods("PUT")
	// Delete a specific doctor
	v1.HandleFunc("/doctors/{id:[0-9]+}", app.requirePermissions("doctors:write", app.deleteDoctorHandler)).Methods("DELETE")

	// Create a new patient
	v1.HandleFunc("/patients", app.createPatientHandler).Methods("POST")
	// Get a doctors list by pagination and filters
	v1.HandleFunc("/patients", app.requirePermissions("patients:read", app.SearchPatientHandler)).Methods("GET")
	// Get a specific patient
	v1.HandleFunc("/patients/{id:[0-9]+}", app.requirePermissions("patients:read", app.getPatientHandler)).Methods("GET")
	// Update a specific patient
	v1.HandleFunc("/patients/{id:[0-9]+}", app.requirePermissions("patients:write", app.updatePatientHandler)).Methods("PUT")
	// Delete a specific patient
	v1.HandleFunc("/patients/{id:[0-9]+}", app.requirePermissions("patients:write", app.deletePatientHandler)).Methods("DELETE")

	// Create a new clinic
	v1.HandleFunc("/clinics", app.createClinicHandler).Methods("POST")
	// Get a specific clinic
	v1.HandleFunc("/clinics/{id:[0-9]+}", app.requirePermissions("clinics:read", app.getClinicHandler)).Methods("GET")
	// Get a clinics list by pagination and filters
	v1.HandleFunc("/clinics", app.requirePermissions("clinics:read", app.searchClinicHandler)).Methods("GET")
	// Update a specific clinic
	v1.HandleFunc("/clinics/{id:[0-9]+}", app.requirePermissions("clinics:write", app.updateClinicHandler)).Methods("PUT")
	// Delete a specific clinic
	v1.HandleFunc("/clinics/{id:[0-9]+}", app.requirePermissions("clinics:write", app.deleteClinicHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	// disabled because not needed
	// users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	log.Printf("Starting server on %s\n", app.config.port)

	err := http.ListenAndServe(app.config.port, r)

	log.Fatal("ListenAndServe Err: ", err)

	return app.authenticate(r)
}
