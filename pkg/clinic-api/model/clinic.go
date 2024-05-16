package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ClinicModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}


func (m ClinicModel) Insert(clinic *Clinic) error {
	query := `
		INSERT INTO clinics (name, city, address) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{clinic.Name, clinic.City, clinic.Address}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&clinic.Id,
		&clinic.CreatedAt,
		&clinic.UpdatedAt,
	)
}


func (m ClinicModel) GetAll(name, city string, filters Filters) ([]*Clinic, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, name, city, address
		FROM clinics
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (LOWER(city) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, city, filters.limit(), filters.offset()}

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
	var clinics []*Clinic

	for rows.Next() {
		var clinic Clinic
		err := rows.Scan(&totalRecords, &clinic.Id, &clinic.CreatedAt, &clinic.UpdatedAt, &clinic.Name, &clinic.City, &clinic.Address)
		if err != nil {
			return nil, Metadata{}, err
		}
		clinics = append(clinics, &clinic)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return clinics, metadata, nil
}


func (m ClinicModel) Get(id int) (*Clinic, error) {
	query := `
		SELECT id, created_at, updated_at, name, city, address
		FROM clinics
		WHERE id = $1
		`
	var clinic Clinic
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&clinic.Id,
		&clinic.CreatedAt,
		&clinic.UpdatedAt,
		&clinic.Name,
		&clinic.City,
		&clinic.Address,
	)

	if err != nil {
		return nil, err
	}
	return &clinic, nil
}


func (m ClinicModel) Update(clinic *Clinic) error {
	query := `
		UPDATE clinics
		SET name = $1, city = $2, address = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{
		clinic.Name,
		clinic.City,
		clinic.Address,
		clinic.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&clinic.UpdatedAt)
}


func (m ClinicModel) Delete(id int) error {
	query := `
		DELETE FROM clinics
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
