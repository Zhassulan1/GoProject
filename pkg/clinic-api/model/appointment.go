package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Appointment struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	PatientId  string `json:"patientId"`
	DoctorId   string `json:"doctorId"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Status     string `json:"status"`
}

type AppointmentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m AppointmentModel) Insert(appointment *Appointment) error {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, date, start_time, end_time, status) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{appointment.PatientId, appointment.DoctorId, appointment.Date, appointment.StartTime, appointment.EndTime, appointment.Status}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt)
}

func (m AppointmentModel) Get(id int) (*Appointment, error) {
	query := `
		SELECT id, created_at, updated_at, patient_id, doctor_id, date, start_time, end_time, status
		FROM appointments
		WHERE id = $1
		`
	var appointment Appointment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.PatientId, &appointment.DoctorId, &appointment.Date, &appointment.StartTime, &appointment.EndTime, &appointment.Status)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (m AppointmentModel) Update(appointment *Appointment) error {
	query := `
		UPDATE appointments
		SET patient_id = $1, doctor_id = $2, date = $3, start_time = $4, end_time = $5, status = $6
		WHERE id = $7
		RETURNING updated_at
		`
	args := []interface{}{appointment.PatientId, appointment.DoctorId, appointment.Date, appointment.StartTime, appointment.EndTime, appointment.Status, appointment.Id}
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
