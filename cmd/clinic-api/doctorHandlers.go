package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/validator"

	"github.com/gorilla/mux"
)

func (app *application) createDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Specialty string `json:"specialty"`
		ClinicID  int    `json:"clinic_id"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//* USER PART STARTS
	user := &model.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if
	// any of the checks fail.
	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get an ErrDuplicateEmail error, use the v.AddError() method to manually add
		// a message to the validator instance, and then call our failedValidationResponse
		// helper().
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")

			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "doctors:write", "clinics:read", "patients:read", "appointments:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	//* USER PART ENDS

	// doctor continuation

	doctor := &model.Doctor{
		Name:      input.Name,
		Specialty: input.Specialty,
		ClinicID:  input.ClinicID,
		UserID:    user.ID,
	}

	err = app.models.Doctors.Insert(doctor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// user part
	res.Token = &token.Plaintext
	res.User = user

	app.writeJSON(w, http.StatusCreated, envelope{"doctor": doctor, "user": res}, nil)
}

// ???????????????????????
// ???????????????????????
// ???????????????????????
// ???????????????????????

func (app *application) SearchDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Specialty string `json:"specialty"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Specialty = app.readStrings(qs, "specialty", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "specialty", "created_at", "updated_at",
		// descending sort values
		"-id", "-name", "-specialty", "-created_at", "-updated_at",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	doctors, metadata, err := app.models.Doctors.GetAll(input.Name, input.Specialty, input.Filters)
	if err != nil {
		fmt.Println("We are in search doctor handler", "\nname: ", input.Name, "\nspec:", input.Specialty, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"doctors": doctors, "metadata": metadata}, nil)
}

// ???????????????????????
// ???????????????????????
// ???????????????????????
// ???????????????????????

func (app *application) getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"doctor": doctor}, nil)
	// app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) updateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	fmt.Println(param)
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid doctor ID")
		fmt.Print("\nError: ", err)
		fmt.Print("\nid: ", id, "\n\n\n")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name      *string `json:"name"`
		Specialty *string `json:"specialty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		doctor.Name = *input.Name
	}

	if input.Specialty != nil {
		doctor.Specialty = *input.Specialty
	}

	err = app.models.Doctors.Update(doctor)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"doctor": doctor}, nil)
}

func (app *application) deleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		fmt.Print("Error: ", err, ": ", id)
		return
	}

	err = app.models.Doctors.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
