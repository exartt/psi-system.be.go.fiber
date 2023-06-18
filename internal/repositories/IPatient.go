package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/person"
)

type PatientRepository interface {
	GetByID(ID uint) (*person.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) GetByID(ID uint) (*person.Patient, error) {
	var patient person.Patient
	err := r.db.First(&patient, "id = ?", ID).Error
	return &patient, err
}
