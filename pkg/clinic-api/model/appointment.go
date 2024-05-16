package model

import (
	"context"
	"time"
	"fmt"
)

func (m AppointmentModel) Insert(appointment *Appointment) error {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, date, start_time, end_time, status) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{
		appointment.PatientId,
		appointment.DoctorId,
		appointment.Date,
		appointment.StartTime,
		appointment.EndTime,
		appointment.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&appointment.Id,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
}
func (m AppointmentModel) GetAll(patientId, doctorId, date, status string, filters Filters) ([]*Appointment, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id,  created_at, updated_at, patient_id, doctor_id, date, start_time, end_time, status
		FROM appointments
		WHERE (patient_id::TEXT ILIKE $1 OR $1 = '')
		AND (doctor_id::TEXT ILIKE $2 OR $2 = '')
		AND (date::TEXT ILIKE $3 OR $3 = '')
		AND (status ILIKE $4 OR $4 = '')
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{patientId, doctorId, date, status, filters.limit(), filters.offset()}

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
	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(&totalRecords, &appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.PatientId, &appointment.DoctorId, &appointment.Date, &appointment.StartTime, &appointment.EndTime, &appointment.Status)
		if err != nil {
			return nil, Metadata{}, err
		}
		appointments = append(appointments, &appointment)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return appointments, metadata, nil
}

func (m AppointmentModel) Get(id int) (*Appointment, error) {
	query := `
		SELECT id,
			created_at,
			updated_at,
			patient_id,
			doctor_id,
			date,
			start_time,
			end_time,
			status
		FROM appointments
		WHERE id = $1
	`
	var appointment Appointment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&appointment.Id,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
		&appointment.PatientId,
		&appointment.DoctorId,
		&appointment.Date,
		&appointment.StartTime,
		&appointment.EndTime,
		&appointment.Status,
	)

	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (m AppointmentModel) Update(appointment *Appointment) error {
	query := `
		UPDATE appointments
		SET patient_id = $1,
			doctor_id  = $2,
			date       = $3,
			start_time = $4,
			end_time   = $5,
			status     = $6
		WHERE id = $7
		RETURNING updated_at
		`
	args := []interface{}{
		appointment.PatientId,
		appointment.DoctorId,
		appointment.Date,
		appointment.StartTime,
		appointment.EndTime,
		appointment.Status,
		appointment.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&appointment.UpdatedAt)
}

func (m AppointmentModel) Delete(id int) error {
	query := `
		DELETE FROM appointments
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
