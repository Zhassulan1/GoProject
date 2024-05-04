package filler

import (
	model "github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
	_ "github.com/lib/pq"
)

func PopulateDatabase(models model.Models) error {
	for _, doctor := range doctors {
		models.Doctors.Insert(&doctor)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var doctors = []model.Doctor{
	{Name: "Doctor Ivanov", Specialty: "Therapy"},
	{Name: "Doctor Petrov", Specialty: "Surgery"},
	{Name: "Doctor Sidorova", Specialty: "Pediatrics"},
	{Name: "Doctor Smirnov", Specialty: "Cardiology"},
	{Name: "Doctor Kozlov", Specialty: "Orthopedics"},
	{Name: "Doctor Morozova", Specialty: "Gynecology"},
	{Name: "Doctor Nikitin", Specialty: "Neurology"},
	{Name: "Doctor Fedorov", Specialty: "Ophthalmology"},
	{Name: "Doctor Alexeeva", Specialty: "Endocrinology"},
	{Name: "Doctor Grigoryev", Specialty: "Otolaryngology"},
}