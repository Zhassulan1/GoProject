package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Patient struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
}

type PatientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m PatientModel) Insert(patient *Patient) error {
	// Insert a new patient into the database.
	query := `
		INSERT INTO patients (name, birthdate, gender) 
		VALUES ($1, $2, $3) 
		RETURNING id
		`
	args := []interface{}{patient.Name, patient.Birthdate, patient.Gender}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.ID)
}

func (m PatientModel) Get(id int) (*Patient, error) {
	// Retrieve a specific patient based on its ID.
	query := `
		SELECT id, name, birthdate, gender
		FROM patients
		WHERE id = $1
		`
	var patient Patient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&patient.ID, &patient.Name, &patient.Birthdate, &patient.Gender)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (m PatientModel) Update(patient *Patient) error {
	// Update a specific patient in the database.
	query := `
		UPDATE patients
		SET name = $1, birthdate = $2, gender = $3
		WHERE id = $4
		`
	args := []interface{}{patient.Name, patient.Birthdate, patient.Gender, patient.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m PatientModel) Delete(id int) error {
	// Delete a specific patient from the database.
	query := `
		DELETE FROM patients
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}