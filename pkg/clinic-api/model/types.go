//! Write in this file your types

package model

import (
	"database/sql"
	"log"
	"time"
)

type Doctor struct {
	id         int64     `json:"id"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"updated_at"`
	name       string    `json:"name"`
	specialty  string    `json:"specialty"`
}

type DoctorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
