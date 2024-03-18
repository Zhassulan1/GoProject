package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DoctorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m DoctorModel) Insert(doctor *Doctor) error {
	query := `
		INSERT INTO doctors (name, specialty) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{doctor.Name, doctor.Specialty}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&doctor.Id,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)
}

func (m DoctorModel) Get(id int) (*Doctor, error) {
	query := `
		SELECT id, created_at, updated_at, name, specialty
		FROM doctors
		WHERE id = $1
		`
	var doctor Doctor
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&doctor.Id,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
		&doctor.Name,
		&doctor.Specialty,
	)

	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (m DoctorModel) Update(doctor *Doctor) error {
	query := `
		UPDATE doctors
		SET name = $1, specialty = $2
		WHERE id = $3
		RETURNING updated_at
		`
	args := []interface{}{
		doctor.Name,
		doctor.Specialty,
		doctor.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&doctor.UpdatedAt)
}

func (m DoctorModel) Delete(id int) error {
	query := `
		DELETE FROM doctors
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
