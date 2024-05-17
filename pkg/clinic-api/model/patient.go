package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PatientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m PatientModel) Insert(patient *Patient) error {
	query := `
		INSERT INTO patients (name, birthdate, gender, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at, user_id
		`

	args := []interface{}{
		patient.Name,
		patient.Birthdate,
		patient.Gender,
		patient.UserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&patient.Id,
		&patient.CreatedAt,
		&patient.UpdatedAt,
		&patient.UserID,
	)
}

func (m PatientModel) GetAll(name, gender string, filters Filters) ([]*Patient, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, name, birthdate, gender
		FROM patients
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (gender = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, gender, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0
	var patients []*Patient
	for rows.Next() {
		var patient Patient
		err := rows.Scan(&totalRecords, &patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.Name, &patient.Birthdate, &patient.Gender)
		if err != nil {
			return nil, Metadata{}, err
		}
		patients = append(patients, &patient)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return patients, metadata, nil
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
	err := row.Scan(
		&patient.Id,
		&patient.CreatedAt,
		&patient.UpdatedAt,
		&patient.Name,
		&patient.Birthdate,
		&patient.Gender,
	)

	if err != nil {
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
	args := []interface{}{
		patient.Name,
		patient.Birthdate,
		patient.Gender,
		patient.Id,
	}

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
