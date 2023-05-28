package services

import (
	"errors"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"psi-system.be.go.fiber/internal/repositories"
)

var Logger = logrus.New()

type AppointmentService interface {
	Save(request *appointment.Appointment) error
}

func NewAppointmentService(repo repositories.AppointmentRepository) AppointmentService {
	return &appointmentService{repo: repo}
}

type appointmentService struct {
	repo repositories.AppointmentRepository
}

func (s *appointmentService) Save(request *appointment.Appointment) error {
	if err := s.checkConflict(request); err != nil {
		return err
	}

	err := s.repo.Create(request)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "Save",
		}).Error("Database error: ", err)
		return err
	}

	Logger.WithFields(logrus.Fields{
		"action": "Save",
	}).Info("Appointment saved successfully")

	return nil
}

func (s *appointmentService) checkConflict(appointment *appointment.Appointment) error {
	existingAppointments, err := s.repo.CheckDateByTimeRange(appointment.Start, appointment.End)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "checkConflict",
		}).Error("Error retrieving appointments: ", err)
		return err
	}
	if len(existingAppointments) > 0 {
		Logger.WithFields(logrus.Fields{
			"action": "checkConflict",
		}).Error("Appointment conflict detected")
		return errors.New("Appointment conflict detected")
	}

	return nil
}

func (s *appointmentService) Update(id uint, request *appointment.Appointment) error {
	if err := s.checkConflict(request); err != nil {
		return err
	}

	err := s.repo.Update(id, request)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "Update",
		}).Error("Database error: ", err)
		return err
	}

	Logger.WithFields(logrus.Fields{
		"action": "Update",
	}).Info("Appointment updated successfully")

	return nil
}

func (s *appointmentService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "Delete",
		}).Error("Database error: ", err)
		return err
	}

	Logger.WithFields(logrus.Fields{
		"action": "Delete",
	}).Info("Appointment deleted successfully")

	return nil
}

func (s *appointmentService) GetAll() ([]*appointment.Appointment, error) {
	appointments, err := s.repo.GetAll()
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "GetAll",
		}).Error("Database error: ", err)
		return nil, err
	}

	Logger.WithFields(logrus.Fields{
		"action": "GetAll",
	}).Info("All appointments retrieved successfully")

	return appointments, nil
}

func (s *appointmentService) UpdateStatusAppointment(ID uint, status enums.StatusAgendamento) error {
	err := s.repo.UpdateStatusAppointment(ID, status)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"action": "UpdateStatusAppointment",
		}).Error("Database error: ", err)
		return err
	}

	Logger.WithFields(logrus.Fields{
		"action": "UpdateStatusAppointment",
	}).Info("Appointment status updated successfully")

	return nil
}
