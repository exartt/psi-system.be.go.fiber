package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/internal/domain/services"
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
