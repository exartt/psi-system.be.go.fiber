package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/services"
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
	options, err := h.Service.GetPatientsOptions()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Error retrieving patients")
	}
	return c.JSON(options)
}
