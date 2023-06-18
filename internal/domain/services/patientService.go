package services

import (
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/repositories"
)

type PatientService interface {
	GetByID(ID uint) (*person.Patient, error)
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
