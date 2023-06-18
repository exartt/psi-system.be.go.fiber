package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"psi-system.be.go.fiber/internal/domain/services"
	"strconv"
	"time"
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

	appointment.CalendarID = os.Getenv("CALENDAR_ID")

	if err := createGoogleCalendarEvent(&appointment, h.Logger); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointmentGoogleCalendar",
		}).Error("Failed to create the google calendar event: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Service.Save(&appointment); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointment",
		}).Error("Error saving appointment: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving appointment"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Appointment created with success!"})
}

func (h *AppointmentHandler) GetAppointmentsByYear(c *fiber.Ctx) error {
	yearStr := c.Query("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "GetAppointmentsByYear",
		}).Error("Invalid year format: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid year format"})
	}

	appointments, err := h.Service.GetByYear(year)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "GetAppointmentsByYear",
		}).Error("Error retrieving appointments: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving appointments"})
	}

	return c.JSON(appointments)
}

func (h *AppointmentHandler) UpdateAppointment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.respondError(c, "UpdateAppointment", "Invalid ID", fiber.StatusBadRequest, err)
	}

	currentAppointment, err := h.Service.GetByID(uint(id))

	if err != nil {
		return h.respondError(c, "UpdateAppointment", "Error fetching appointment", fiber.StatusInternalServerError, err)
	}

	updates, err := h.parseUpdateData(c)
	if err != nil {
		return h.respondError(c, "UpdateAppointment", "Invalid request body", fiber.StatusBadRequest, err)
	}

	if h.applyUpdates(currentAppointment, updates) {
		if err := h.Service.Update(uint(id), currentAppointment); err != nil {
			return h.respondError(c, "UpdateAppointment", "Error updating appointment", fiber.StatusInternalServerError, err)
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Appointment status updated successfully"})
}

func (h *AppointmentHandler) CancelAppointment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return h.respondError(c, "CancelAppointment", "Invalid ID", fiber.StatusBadRequest, err)
	}

	currentAppointment, err := h.Service.GetByID(uint(id))
	if err != nil {
		return h.respondError(c, "CancelAppointment", "Error fetching appointment", fiber.StatusInternalServerError, err)
	}

	// Definir o status como "Cancelado"
	currentAppointment.Status = enums.Cancelado

	// Atualizar o compromisso no banco de dados
	if err := h.Service.Update(uint(id), currentAppointment); err != nil {
		return h.respondError(c, "CancelAppointment", "Error updating appointment", fiber.StatusInternalServerError, err)
	}

	// Enviar resposta bem-sucedida
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Appointment canceled successfully"})
}

func (h *AppointmentHandler) parseUpdateData(c *fiber.Ctx) (*appointment.Appointment, error) {
	var updates appointment.Appointment
	err := c.BodyParser(&updates)
	return &updates, err
}

func (h *AppointmentHandler) applyUpdates(appointment *appointment.Appointment, updates *appointment.Appointment) bool {
	hasChanges := false

	if !updates.Start.IsZero() && !updates.Start.Equal(appointment.Start) {
		appointment.Start = updates.Start
		hasChanges = true
	}

	if !updates.End.IsZero() && !updates.End.Equal(appointment.End) {
		appointment.End = updates.End
		hasChanges = true
	}

	if updates.PsychologistID != 0 && updates.PsychologistID != appointment.PsychologistID {
		appointment.PsychologistID = updates.PsychologistID
		hasChanges = true
	}

	if updates.PatientID != 0 && updates.PatientID != appointment.PatientID {
		appointment.PatientID = updates.PatientID
		hasChanges = true
	}

	if updates.Summary != "" && updates.Summary != appointment.Summary {
		appointment.Summary = updates.Summary
		hasChanges = true
	}

	if updates.Description != "" && updates.Description != appointment.Description {
		appointment.Description = updates.Description
		hasChanges = true
	}

	if updates.Notify != appointment.Notify {
		appointment.Notify = updates.Notify
		hasChanges = true
	}

	return hasChanges
}

func (h *AppointmentHandler) respondError(c *fiber.Ctx, action, message string, statusCode int, err error) error {
	h.Logger.WithFields(logrus.Fields{
		"action": action,
	}).Error(message+": ", err)
	return c.Status(statusCode).JSON(fiber.Map{"error": message})
}

func createGoogleCalendarEvent(appointment *appointment.Appointment, logger *logrus.Logger) error {
	calendarEvent := map[string]interface{}{
		"summary":     "Consulta",
		"description": appointment.Description,
		"location":    "Clínica",
		"start": map[string]string{
			"dateTime": appointment.Start.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo", // Brasil (GMT-3)
		},
		"end": map[string]string{
			"dateTime": appointment.End.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo", // Brasil (GMT-3)
		},
	}

	jsonData, err := json.Marshal(calendarEvent)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"action": "createGoogleCalendarEvent",
		}).Error("Failed to marshal JSON: ", err)
		return fmt.Errorf("Failed to marshal JSON")
	}

	resp, err := http.Post(enums.CREATE.String(""), "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.WithFields(logrus.Fields{
			"action": "createGoogleCalendarEvent",
		}).Error("Failed to create calendar event: ", err)
		return fmt.Errorf("Failed to create calendar event")
	}

	return nil
}
