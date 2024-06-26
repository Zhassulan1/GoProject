package model

import "time"

type Patient struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
	UserID    int64  `json:"user_id"`
}

type Doctor struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Specialty string `json:"specialty"`
	ClinicID  int    `json:"clinic_id"`
	UserID    int64  `json:"user_id"`
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type DoctorPagination struct {
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Name         string `json:"name"`
	Specialty    string `json:"specialty"`
	CreatedBegin string `json:"startTime"`
	CreatedEnd   string `json:"endTime"`
	OrderBy      string `json:"orederBy"`
	Order        uint8  `json:"order"`
}

type Appointment struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	PatientId string `json:"patientId"`
	DoctorId  string `json:"doctorId"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    string `json:"status"`
}

type Clinic struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	City      string `json:"city"`
}
