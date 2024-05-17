package model

import (
	"context"
	"database/sql"
	"fmt"
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
		INSERT INTO doctors (name, specialty, clinic_id, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at, user_id
		`
	args := []interface{}{
		doctor.Name,
		doctor.Specialty,
		doctor.ClinicID,
		doctor.UserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&doctor.Id,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
		&doctor.UserID,
	)
}

func (m DoctorModel) GetAll(name, specialty string, filters Filters) ([]*Doctor, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, name, specialty, clinic_id
		FROM doctors
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (LOWER(specialty) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, specialty, filters.limit(), filters.offset()}

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
	var doctors []*Doctor
	for rows.Next() {
		var doctor Doctor
		err := rows.Scan(&totalRecords, &doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.Name, &doctor.Specialty, &doctor.ClinicID)
		if err != nil {
			return nil, Metadata{}, err
		}
		doctors = append(doctors, &doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return doctors, metadata, nil
}

func (m DoctorModel) Get(id int) (*Doctor, error) {
	query := `
		SELECT id, created_at, updated_at, name, specialty, clinic_id
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
		&doctor.ClinicID,
	)

	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (m DoctorModel) Update(doctor *Doctor) error {
	query := `
		UPDATE doctors
		SET name = $1, specialty = $2, clinic_id = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{
		doctor.Name,
		doctor.Specialty,
		doctor.ClinicID,
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

func (m DoctorModel) GetByClinicID(clinicID int) ([]*Doctor, error) {
    query := `SELECT id, name, specialty FROM doctors WHERE clinic_id = $1`

    rows, err := m.DB.Query(query, clinicID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var doctors []*Doctor

    for rows.Next() {
        var doctor Doctor
        err = rows.Scan(&doctor.Id, &doctor.Name, &doctor.Specialty)
        if err != nil {
            return nil, err
        }
        doctors = append(doctors, &doctor)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return doctors, nil
}