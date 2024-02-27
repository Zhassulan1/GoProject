package main

// import "database/sql"

// type Strorage interface {
// 	DoctorInfo(*Doctor) error
// 	CreateDoctor(*Doctor) error
// }

// type PostgresStore struct {
// 	db *sql.DB
// }

// func NewPostgresStore() (*PostgresStore, error) {
// 	connStr := "user=postgres dbname=clinic_db password=1234 sslmode=disable"
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return &PostgresStore{
// 		db: db,
// 	}, nil
// }

// func (s *PostgresStore) Init() error {
// 	return s.createDoctorTable()
// }

// func (s *PostgresStore) createDoctorTable() error {
// 	query := `create table if not exists account (
// 		id serial primary key,
// 		first_name varchar(100),
// 		last_name varchar(100),
// 		number serial,
// 		encrypted_password varchar(100),
// 		balance serial,
// 		created_at timestamp
// 	)`

// 	_, err := s.db.Exec(query)
// 	return err
// }
