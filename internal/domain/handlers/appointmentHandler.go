package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"psi-system.be.go.fiber/internal/domain/services"
)

type AppointmentHandler struct {
	Service services.AppointmentService
	Logger  *logrus.Logger
}

func NewAppointmentHandler(service services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		Service: service,
	}
}

func (h *AppointmentHandler) CreateAppointment(c *fiber.Ctx) error {
	var appointment appointment.Appointment
	if err := c.BodyParser(&appointment); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointment",
		}).Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	if err := h.Service.Save(&appointment); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointment",
		}).Error("Error saving appointment: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving appointment"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Appointment created with success!"})
}
