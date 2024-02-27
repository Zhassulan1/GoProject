package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	// Menus       MenuModel
	// Restaurants RestaurantModel

	Doctors DoctorModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		// Menus: MenuModel{
		// 	DB:       db,
		// 	InfoLog:  infoLog,
		// 	ErrorLog: errorLog,
		// },
		Doctors: DoctorModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
