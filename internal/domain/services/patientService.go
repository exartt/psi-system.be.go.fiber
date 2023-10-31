package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/repositories"
	"strconv"
	"time"
)

type Option struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	SessionPrice string `json:"session_price"`
}

type PatientService interface {
	GetByID(ID uint) (*person.Patient, error)
	GetPatientsOptions(psychologistID uint) ([]Option, error)
	Create(patient person.DTO, psychologistID uint) (*person.Patient, error)
	Update(patient person.DTO, psychologistID uint) (*person.Patient, error)
	Delete(ID uint) error
	GetPersonPatient(psychologistID uint) ([]person.PersonPatient, error)
	GetPatient(psychologistID uint, patientID uint) (person.DTO, error)
	DeactivatePatient(ID uint) error
}

type patientService struct {
	repo       repositories.PatientRepository
	personRepo repositories.IPersonRepository
}

func NewPatientService(repo repositories.PatientRepository, personRepo repositories.IPersonRepository) PatientService {
	return &patientService{repo: repo, personRepo: personRepo}
}

func (s *patientService) GetByID(ID uint) (*person.Patient, error) {
	patient, err := s.repo.GetByID(ID)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "GetPatientByID",
		}).Error("Database error: ", err)
		return nil, err
	}

	Logger.WithFields(logrus.Fields{
		"action": "GetPatientByID",
	}).Info("Patient retrieved successfully")

	return patient, nil
}

func (s *patientService) GetPatientsOptions(psychologistID uint) ([]Option, error) {
	patients, err := s.repo.GetPatientsWithPersonName(psychologistID)
	if err != nil {
		return nil, err
	}

	var options []Option
	for _, patient := range patients {
		options = append(options, Option{
			Value:        fmt.Sprint(patient.ID),
			Label:        patient.PersonName,
			SessionPrice: patient.SessionPrice,
		})
	}
	return options, nil
}

func (s *patientService) GetPersonPatient(psychologistID uint) ([]person.PersonPatient, error) {
	personPatient, err := s.repo.GetPersonPatient(psychologistID)
	if err != nil {
		return nil, err
	}
	if personPatient == nil {
		personPatient = []person.PersonPatient{}
	}
	return personPatient, nil
}

func (s *patientService) Create(patient person.DTO, psychologistID uint) (*person.Patient, error) {
	personDto := &person.Person{
		Name:      patient.Name,
		Email:     patient.Email,
		CellPhone: patient.CellPhone,
		Phone:     patient.Phone,
		ZipCode:   patient.ZipCode,
		Address:   patient.Address,
		IsActive:  patient.IsActive,
		CPF:       patient.CPF,
		RG:        patient.RG,
	}

	if err := s.personRepo.CreatePerson(personDto); err != nil {
		print(err)
		return nil, err
	}

	sessionPrice, err := strconv.ParseFloat(patient.SessionPrice, 32)
	if err != nil {
		sessionPrice = 140
	}

	patientDto := &person.Patient{
		IsPlan:         patient.IsPlan,
		SessionPrice:   sessionPrice,
		ConversionType: patient.ConversionType,
		PsychologistID: psychologistID,
	}

	patientDto.PersonID = personDto.ID

	if err := s.repo.Create(patientDto); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *patientService) GetPatient(psychologistID uint, patientID uint) (person.DTO, error) {
	dto, err := s.repo.GetPatient(psychologistID, patientID)
	if err != nil {
		return person.DTO{}, err
	}
	return dto, nil
}

func (s *patientService) Update(patient person.DTO, psychologistID uint) (*person.Patient, error) {
	personDto := &person.Person{
		ID:        patient.PersonId,
		Name:      patient.Name,
		Email:     patient.Email,
		CellPhone: patient.CellPhone,
		Phone:     patient.Phone,
		ZipCode:   patient.ZipCode,
		Address:   patient.Address,
		IsActive:  patient.IsActive,
		CPF:       patient.CPF,
		RG:        patient.RG,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.personRepo.Update(personDto); err != nil {
		print(err)
		return nil, err
	}

	sessionPrice, err := strconv.ParseFloat(patient.SessionPrice, 32)

	patientDto := &person.Patient{
		ID:             patient.ID,
		IsPlan:         patient.IsPlan,
		SessionPrice:   sessionPrice,
		ConversionType: patient.ConversionType,
		PsychologistID: psychologistID,
	}

	patientDto.PersonID = personDto.ID
	err = s.repo.Update(patientDto)
	return patientDto, err
}

func (s *patientService) Delete(ID uint) error {
	return s.repo.Delete(ID)
}

func (s *patientService) DeactivatePatient(ID uint) error {
	return s.repo.DeactivatePatient(ID)
}
