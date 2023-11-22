package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/domain/utils"
	"strconv"
	"strings"
)

type PatientHandler struct {
	Service services.PatientService
	Logger  *logrus.Logger
}

func NewPatientHandler(service services.PatientService) *PatientHandler {
	return &PatientHandler{
		Service: service,
	}
}

func (h *PatientHandler) GetPatientsOptions(c *fiber.Ctx) error {
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	options, err := h.Service.GetPatientsOptions(psychologistID)
	if err != nil || len(options) == 0 {
		fmt.Println(err)
		return c.Status(500).SendString("Não há pacientes cadastrados. Para cadastrar um novo paciente, acesse a aba 'Pacientes' no menu lateral e clique em 'Novo Paciente'")
	}
	return c.JSON(options)
}

func (h *PatientHandler) CreatePatient(c *fiber.Ctx) error {
	var patient person.DTO
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err := c.BodyParser(&patient); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Não foi possível analisar o corpo da requisição"})
	}

	newPatient, err := h.Service.Create(patient, psychologistID)
	if err != nil {
		if strings.Contains(err.Error(), "people_email_key") {
			return c.Status(400).JSON(fiber.Map{"error": "Email já cadastrado"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Erro interno ao criar paciente"})
	}

	return c.JSON(newPatient)
}
func (h *PatientHandler) UpdatePatient(c *fiber.Ctx) error {
	var patient person.DTO
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err := c.BodyParser(&patient); err != nil {
		return c.Status(400).SendString("Bad Request")
	}

	updatedPatient, err := h.Service.Update(patient, psychologistID)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(updatedPatient)
}

func (h *PatientHandler) DeletePatient(c *fiber.Ctx) error {
	ID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	if err := h.Service.Delete(uint(ID)); err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString("Patient deleted")
}

func (h *PatientHandler) GetPersonPatient(c *fiber.Ctx) error {
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	personPatient, err := h.Service.GetPersonPatient(psychologistID)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(personPatient)
}

func (h *PatientHandler) DeactivatePatient(c *fiber.Ctx) error {
	ID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	if err := h.Service.DeactivatePatient(uint(ID)); err != nil {
		if err.Error() == "Patient not found" {
			return c.Status(404).SendString("Patient not found")
		}
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.SendString("Patient deactivated")
}

func (h *PatientHandler) GetPatient(c *fiber.Ctx) error {
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err != nil {
		return c.Status(400).SendString("Bad Request: Invalid Psychologist ID")
	}

	patientIDStr := c.Params("id")
	patientID, err := strconv.ParseUint(patientIDStr, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Bad Request: Invalid Patient ID")
	}

	dto, err := h.Service.GetPatient(psychologistID, uint(patientID))
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	return c.JSON(dto)
}

func (h *PatientHandler) GetPatientOption(ctx *fiber.Ctx) error {
	psychologistID, err := utils.GetPsychologistIDFromContext(ctx)
	if err != nil {
		return ctx.Status(400).SendString("Bad Request: Invalid Psychologist ID")
	}

	patientIDStr := ctx.Params("id")
	patientID, err := strconv.ParseUint(patientIDStr, 10, 64)
	if err != nil {
		return ctx.Status(400).SendString("Bad Request: Invalid Patient ID")
	}

	dto, err := h.Service.GetPatientOption(psychologistID, uint(patientID))
	if err != nil {
		return ctx.Status(500).SendString("Internal Server Error")
	}

	return ctx.JSON(dto)
}
