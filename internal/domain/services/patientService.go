package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/repositories"
)

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type PatientService interface {
	GetByID(ID uint) (*person.Patient, error)
	GetPatientsOptions() ([]Option, error)
}

type patientService struct {
	repo repositories.PatientRepository
}

func NewPatientService(repo repositories.PatientRepository) PatientService {
	return &patientService{repo: repo}
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

func (s *patientService) GetPatientsOptions() ([]Option, error) {
	patients, err := s.repo.GetPatientsWithPersonName()
	if err != nil {
		return nil, err
	}

	var options []Option
	for _, patient := range patients {
		options = append(options, Option{
			Value: fmt.Sprint(patient.ID),
			Label: patient.PersonName,
		})
	}
	return options, nil
}
