package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/domain/utils"
)

type GoogleConsumerHandler struct {
	Service services.GoogleConsumerService
	Logger  *logrus.Logger
}

func NewGoogleConsumerHandler(service services.GoogleConsumerService) *GoogleConsumerHandler {
	return &GoogleConsumerHandler{
		Service: service,
	}
}

func (h *GoogleConsumerHandler) RequestGoogleAuth(c *fiber.Ctx) error {
	url, err := h.Service.RequestGoogleAuth()
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "RequestGoogleAuth",
		}).Error("Error requesting google auth: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error requesting google auth"})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"auth_url": url})
}

func (h *GoogleConsumerHandler) RequestGoogleAuthAuthorized(c *fiber.Ctx) error {
	token, err := h.Service.RequestGoogleAuthAuthorized()
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"action": "RequestGoogleAuthAuthorized",
		}).Error("Error requesting Google authorized auth: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error requesting Google authorized auth"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"access_token": token})
}

func (h *GoogleConsumerHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	accessToken := c.FormValue("access_token")
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro inesperado"})
	}

	expirationTime := utils.GetExpirationDate()
	err = h.Service.Store(accessToken, expirationTime, psychologistID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not store token"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"access_token": accessToken})
}
