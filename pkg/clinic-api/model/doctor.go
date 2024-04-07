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

// ???????????????????????
// ???????????????????????
// ???????????????????????
// ???????????????????????

func (m DoctorModel) GetAll(name, specialty string, filters Filters) ([]*Doctor, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, name, specialty
		FROM doctors
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (LOWER(specialty) = LOWER($2) OR $2 = '')
		-- AND (nutrition_value >= $2 OR $2 = 0)
		-- AND (nutrition_value <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{name, specialty, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var doctors []*Doctor
	for rows.Next() {
		var doctor Doctor
		err := rows.Scan(&totalRecords, &doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.Name, &doctor.Specialty)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		doctors = append(doctors, &doctor)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return doctors, metadata, nil
}

// ???????????????????????
// ???????????????????????
// ???????????????????????
// ???????????????????????

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
