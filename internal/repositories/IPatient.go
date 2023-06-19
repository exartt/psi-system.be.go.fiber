package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/person"
)

type PatientRepository interface {
	GetByID(ID uint) (*person.Patient, error)
	GetAll() ([]*person.Patient, error)
	GetPatientsWithPersonName() ([]PatientResult, error)
}

type patientRepository struct {
	db *gorm.DB
}

type PatientResult struct {
	ID         uint
	PersonName string
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) GetByID(ID uint) (*person.Patient, error) {
	var patient person.Patient
	err := r.db.First(&patient, "id = ?", ID).Error
	return &patient, err
}

func (r *patientRepository) GetAll() ([]*person.Patient, error) {
	var patient []*person.Patient
	err := r.db.Find(&patient).Error
	return patient, err
}

func (r *patientRepository) GetPatientsWithPersonName() ([]PatientResult, error) {
	var patients []PatientResult

	err := r.db.Table("patients").
		Select("patients.id, people.name as person_name").
		Joins("JOIN people ON people.id = patients.person_id").
		Scan(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}
