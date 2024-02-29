package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Patient struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
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
	query := `
		INSERT INTO patients (name, birthdate, gender) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{patient.Name, patient.Birthdate, patient.Gender}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt)
}

func (m PatientModel) Get(id int) (*Patient, error) {

	query := `
		SELECT id, created_at, updated_at, name, birthdate, gender
		FROM patients
		WHERE id = $1
		`
	var patient Patient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)

	log.Print(row.Err(), "\n", row.Scan(), "\n", row)

	err := row.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.Name, &patient.Birthdate, &patient.Gender)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return &patient, nil
}

func (m PatientModel) Update(patient *Patient) error {
	query := `
		UPDATE patients
		SET name = $1, birthdate = $2, gender = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{patient.Name, patient.Birthdate, patient.Gender, patient.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.UpdatedAt)
}

func (m PatientModel) Delete(id int) error {
	query := `
		DELETE FROM patients
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
