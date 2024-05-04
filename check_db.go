package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// DSN string anatomy
	// Linux/Mac: postgres://username:password@host:port/database?sslmode=disable
	// Windows: postgresql://username:password@host:port/database?sslmode=disable
	db, err := sql.Open("postgres", "postgres://postgres:1234@db:5432/medicalclinic?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to DB!")
}
