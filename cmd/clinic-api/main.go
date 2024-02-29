package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1234@localhost/medicalclinic?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// CLinic Singleton
	// Create a new appointment
	v1.HandleFunc("/appointments", app.createAppointmentHandler).Methods("POST")
	// Get a specific appointment
	v1.HandleFunc("/appointments/{appointmentId:[0-9]+}", app.getAppointmentHandler).Methods("GET")
	// Update a specific appointment
	v1.HandleFunc("/appointments/{appointmentId:[0-9]+}", app.updateAppointmentHandler).Methods("PUT")
	// Delete a specific appointment
	v1.HandleFunc("/appointments/{appointmentId:[0-9]+}", app.deleteAppointmentHandler).Methods("DELETE")

	// Create a new doctor
	v1.HandleFunc("/doctors", app.createDoctorHandler).Methods("POST")
	// Get a specific doctor
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.getDoctorHandler).Methods("GET")
	// Update a specific doctor
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.updateDoctorHandler).Methods("PUT")
	// Delete a specific doctor
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.deleteDoctorHandler).Methods("DELETE")

	// Create a new patient
	v1.HandleFunc("/patients", app.createPatientHandler).Methods("POST")
	// Get a specific patient
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.getPatientHandler).Methods("GET")
	// Update a specific patient
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.updatePatientHandler).Methods("PUT")
	// Delete a specific patient
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.deletePatientHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection is working by executing a test query
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
