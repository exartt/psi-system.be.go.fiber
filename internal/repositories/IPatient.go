package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/person"
)

type PatientRepository interface {
	GetByID(ID uint) (*person.Patient, error)
	GetAll() ([]*person.Patient, error)
	GetPatientsWithPersonName(psychologistID uint) ([]PatientResult, error)
	Create(*person.Patient) error
	Update(*person.Patient) error
	Delete(ID uint) error
	GetPersonPatient(psychologistID uint) ([]person.PersonPatient, error)
	GetPatient(psychologistID uint, patientID uint) (person.DTO, error)
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

func (r *patientRepository) GetPatientsWithPersonName(psychologistID uint) ([]PatientResult, error) {
	var patients []PatientResult

	err := r.db.Table("patients").
		Select("patients.id, people.name as person_name").
		Joins("JOIN people ON people.id = patients.person_id").
		Where("patients.psychologist_id = ?", psychologistID).
		Scan(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) GetPersonPatient(psychologistID uint) ([]person.PersonPatient, error) {
	var patients []person.PersonPatient

	err := r.db.Table("patients").
		Select("patients.id, people.name as Name, people.email as Email").
		Joins("JOIN people ON people.id = patients.person_id").
		Where("patients.psychologist_id = ?", psychologistID).
		Scan(&patients).Error

	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) GetPatient(psychologistID uint, patientID uint) (person.DTO, error) {
	var dto person.DTO

	err := r.db.Table("patients").
		Select("patients.id, people.name as \"Name\", people.id as \"PersonId\", people.email as \"Email\", people.cell_phone as \"CellPhone\", "+
			"people.phone as \"Phone\", people.zip_code as \"ZipCode\", people.address as \"Address\", people.is_active as \"IsActive\","+
			" people.cpf as \"CPF\", people.rg as \"RG\", "+
			"patients.is_plan as \"IsPlan\", patients.session_price as \"SessionPrice\", patients.conversion_type as \"ConversionType\"").
		Joins("JOIN people ON people.id = patients.person_id").
		Where("patients.psychologist_id = ? AND patients.id = ?", psychologistID, patientID).
		Scan(&dto).Error

	if err != nil {
		return person.DTO{}, err
	}

	return dto, nil
}

func (r *patientRepository) Create(patient *person.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) Update(patient *person.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(ID uint) error {
	return r.db.Delete(&person.Patient{}, ID).Error
}
