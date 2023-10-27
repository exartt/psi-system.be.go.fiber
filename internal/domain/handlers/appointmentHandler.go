package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/appointment"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/domain/utils"
	"psi-system.be.go.fiber/internal/infrastructure"
	"psi-system.be.go.fiber/internal/repositories"
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
	var appointmentDTO appointment.DTO

	if err := c.BodyParser(&appointmentDTO); err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointment",
		}).Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	appointment := convertDTOToEntity(appointmentDTO)
	googleRepo := repositories.NewGCalendarRepository(infrastructure.DB)
	psyID, _ := utils.GetPsychologistIDFromContext(c)
	token, err := googleRepo.FindByPsyID(psyID)
	if err != nil {
		return err
	}

	eventID, err := createGoogleCalendarEvent(&appointment, h.Logger, token.Token)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "CreateAppointmentGoogleCalendar",
		}).Error("Failed to create the google calendar event: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	appointment.EventID = eventID
	appointment.PsychologistID = psyID
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

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

	googleRepo := repositories.NewGCalendarRepository(infrastructure.DB)
	psyID, _ := utils.GetPsychologistIDFromContext(c)
	token, err := googleRepo.FindByPsyID(psyID)
	if err != nil {
		return err
	}
	if h.applyUpdates(currentAppointment, updates) {
		if err := updateGoogleCalendarEvent(updates, updates.EventID, h.Logger, token.Token); err != nil {
			return h.respondError(c, "UpdateAppointment", "Error updating Google Calendar event", fiber.StatusInternalServerError, err)
		}

		currentAppointment.Status = enums.Remarcado
		currentAppointment.UpdatedAt = time.Now()
		currentAppointment.PsychologistID = psyID

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

	googleRepo := repositories.NewGCalendarRepository(infrastructure.DB)
	psyID, _ := utils.GetPsychologistIDFromContext(c)
	token, err := googleRepo.FindByPsyID(psyID)
	if err != nil {
		return err
	}

	if err := deleteGoogleCalendarEvent(currentAppointment.EventID, h.Logger, token.Token); err != nil {
		return h.respondError(c, "CancelAppointment", "Error canceling Google Calendar event", fiber.StatusInternalServerError, err)
	}

	currentAppointment.Status = enums.Cancelado
	currentAppointment.UpdatedAt = time.Now()
	currentAppointment.PsychologistID = psyID

	if err := h.Service.Update(uint(id), currentAppointment); err != nil {
		return h.respondError(c, "CancelAppointment", "Error updating appointment", fiber.StatusInternalServerError, err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Appointment canceled successfully"})
}

func (h *AppointmentHandler) parseUpdateData(c *fiber.Ctx) (*appointment.Appointment, error) {
	var updates appointment.DTO
	err := c.BodyParser(&updates)
	appointment := convertDTOToEntity(updates)
	return &appointment, err
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

func deleteGoogleCalendarEvent(eventID string, logger *logrus.Logger, token string) error {
	client := &http.Client{}

	deleteEventURL := enums.DELETE.String(eventID)

	req, err := http.NewRequest(http.MethodDelete, deleteEventURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"action": "deleteGoogleCalendarEvent",
		}).Error("Failed to create request: ", err)
		return fmt.Errorf("Failed to create request")
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil && resp.StatusCode == 410 {
			logger.WithFields(logrus.Fields{
				"action": "deleteGoogleCalendarEvent",
			}).Info("Resource has already been deleted. Ignoring.")
			return nil
		}
		logger.WithFields(logrus.Fields{
			"action": "deleteGoogleCalendarEvent",
		}).Error("Failed to delete event: ", err)
		return fmt.Errorf("Failed to delete event")
	}

	return nil
}

func createGoogleCalendarEvent(appointment *appointment.Appointment, logger *logrus.Logger, token string) (string, error) {
	calendarEvent := map[string]interface{}{
		"summary":     appointment.Summary,
		"description": appointment.Description,
		"location":    "Clínica",
		"start": map[string]string{
			"dateTime": appointment.Start.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo",
		},
		"end": map[string]string{
			"dateTime": appointment.End.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo",
		},
	}
	fmt.Printf("%+v\n", calendarEvent)
	jsonData, err := json.Marshal(calendarEvent)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"action": "createGoogleCalendarEvent",
		}).Error("Failed to marshal JSON: ", err)
		return "", fmt.Errorf("Failed to marshal JSON")
	}

	req, err := http.NewRequest("POST", enums.CREATE.String(""), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		// Log the error
		logger.WithFields(logrus.Fields{
			"action": "createGoogleCalendarEvent",
		}).Error("Failed to create calendar event: ", err)

		if resp != nil {
			resp.Body.Close()
		}

		return "", fmt.Errorf("Failed to create calendar event")
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body")
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		return "", fmt.Errorf("Failed to parse response body")
	}

	eventID, ok := responseMap["eventID"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to extract event ID")
	}

	return eventID, nil
}

func updateGoogleCalendarEvent(appointment *appointment.Appointment, eventID string, logger *logrus.Logger, token string) error {
	calendarEvent := map[string]interface{}{
		"summary":     appointment.Summary,
		"description": appointment.Description,
		"location":    "Clínica",
		"start": map[string]string{
			"dateTime": appointment.Start.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo",
		},
		"end": map[string]string{
			"dateTime": appointment.End.Format(time.RFC3339),
			"timeZone": "America/Sao_Paulo",
		},
	}

	jsonData, err := json.Marshal(calendarEvent)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"action": "updateGoogleCalendarEvent",
		}).Error("Failed to marshal JSON: ", err)
		return fmt.Errorf("Failed to marshal JSON")
	}

	client := &http.Client{}

	updateEventURL := enums.UPDATE.String(eventID)
	print(updateEventURL)
	req, err := http.NewRequest(http.MethodPut, updateEventURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"action": "updateGoogleCalendarEvent",
		}).Error("Failed to create request: ", err)
		return fmt.Errorf("Failed to create request")
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.WithFields(logrus.Fields{
			"action": "updateGoogleCalendarEvent",
		}).Error("Failed to update calendar event: ", err)
		return fmt.Errorf("Failed to update calendar event")
	}

	return nil
}

func convertDTOToEntity(dto appointment.DTO) appointment.Appointment {
	return appointment.Appointment{
		ID:             dto.ID,
		PsychologistID: dto.PsychologistID,
		PatientID:      dto.PatientID,
		Start:          dto.Start,
		End:            dto.End,
		Summary:        dto.Summary,
		Description:    dto.Description,
		Status:         dto.Status,
		EventID:        dto.EventID,
	}
}
