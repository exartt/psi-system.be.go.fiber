package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/domain/services"
)

type PsychologistHandler struct {
	Service services.PsychologistService
	Logger  *logrus.Logger
}

func NewPsychologistHandler(service services.PsychologistService) *PsychologistHandler {
	return &PsychologistHandler{
		Service: service,
	}
}

func (h *PsychologistHandler) CreatePsychologist(c *fiber.Ctx) error {
	var psychologistDTO person.PsychologistDto

	if err := c.BodyParser(&psychologistDTO); err != nil {
		h.Logger.Error("Failed to read body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	if err := h.Service.SavePsychologist(&psychologistDTO); err != nil {
		h.Logger.Error("Error saving psychologist: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving psychologist"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Psychologist created with success!"})
}
