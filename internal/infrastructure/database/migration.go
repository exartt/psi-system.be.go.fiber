package database

import (
	"fmt"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/infrastructure"
)

func Migrate() error {
	err := infrastructure.DB.AutoMigrate(&appointment.Appointment{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate Appointment model: %v", err)
	}

	err = infrastructure.DB.AutoMigrate(&person.Person{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate Person model: %v", err)
	}

	err = infrastructure.DB.AutoMigrate(&person.Psychologist{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate Psychologist model: %v", err)
	}

	err = infrastructure.DB.AutoMigrate(&person.Patient{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate Patient model: %v", err)
	}

	return nil
}
